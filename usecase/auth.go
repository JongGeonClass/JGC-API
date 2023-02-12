package usecase

import (
	"context"

	"github.com/JongGeonClass/JGC-API/config"
	"github.com/JongGeonClass/JGC-API/database"
	"github.com/JongGeonClass/JGC-API/dbmodel"
	"github.com/JongGeonClass/JGC-API/util"
)

// Auth Usecase의 인터페이스입니다.
type AuthUsecase interface {
	SignUp(ctx context.Context, email, nickname, username, password string) (int64, error)
	Login(ctx context.Context, username, password string) (string, error)
}

// Auth Usecase의 구현체입니다.
type AuthUC struct {
	userdb database.UserDatabase
}

// 회원가입합니다.
// 리턴 타입의 int64는 생성된 유저의 id입니다.
// 만약 이미 존재하는 닉네임을 가진 유저라면 -1을 반환합니다.
// 만약 이미 존재하는 아이디(유저네임)을 가진 유저라면 -2를 반환합니다.
func (uc *AuthUC) SignUp(ctx context.Context, email, nickname, username, password string) (int64, error) {
	user := &dbmodel.User{
		Email:    email,
		Nickname: nickname,
		Username: username,
		Password: password,
		Salt:     util.NewUuid(),
	}
	user.Password = util.Encrypt256(user.Password, user.Salt)

	// 트랜잭션 시작
	// 여기서 에러가 발생하면, 자동으로 트랜잭션을 롤백합니다.
	// 에러가 발생하지 않으면, 커밋합니다.
	// 유저가 존재하는지 체크하고, 존재하지 않는다면, 유저를 생성합니다.
	err := uc.userdb.ExecTx(ctx, func(txdb database.UserDatabase) error {
		// 같은 nickname을 가진 아이디가 존재하는지 검사합니다.
		if exist, err := txdb.CheckUserExistsByNickname(ctx, nickname); err != nil {
			return err
		} else if exist {
			user.Id = -1
			return nil
		}

		// 같은 username을 가진 아이디가 존재하는지 검사합니다.
		if exist, err := txdb.CheckUserExistsByUsername(ctx, username); err != nil {
			return err
		} else if exist {
			user.Id = -2
			return nil
		}

		// user를 추가합니다.
		uid, err := txdb.AddUser(ctx, user)
		if err != nil {
			return err
		}
		user.Id = uid
		return nil
	})
	return user.Id, err
}

// 로그인합니다.
// 로그인에 성공한다면, 토큰을 발급합니다.
// 로그인에 실패했을 경우 빈 문자열을 반환합니다.
// TODO: 현재는 항상 토큰을 재발급합니다.
// 하지만 추후에는 유효한 토큰이 있다면, 해당 토큰을 반환하는 형태로 변경하거나, 유연하게 처리하도록 변경해야 합니다.
// 하지만 로그인을 요청한다는 것 자체가 토큰이 없다는 가정이 될 수도 있으므로 반영하지 않을 수도 있습니다.
// 다중 기기에서의 로그인 처리에 따라 어떻게 할지 추후에 고민해보도록 합시다.
func (uc *AuthUC) Login(ctx context.Context, username, password string) (string, error) {
	type Result struct {
		Token string `json:"token"`
	}
	res := &Result{}
	conf := config.Get()

	// 트랜잭션 시작
	err := uc.userdb.ExecTx(ctx, func(txdb database.UserDatabase) error {
		// 해당 유저가 존재하는지 확인합니다.
		exist, err := txdb.CheckUserExistsByUsername(ctx, username)
		if err != nil {
			return err
		}
		if !exist { // 아이디가 없음
			return nil
		}

		// 유저의 패스워드 정보를 위해 유저 정보를 불러옵니다.
		user, err := txdb.GetUserByUsername(ctx, username)
		if err != nil {
			return err
		}
		if user.Password != util.Encrypt256(password, user.Salt) { // 비밀번호 불일치
			return nil
		}

		// 유저를 인증하는 토큰을 생성합니다.
		tok, err := util.CreateUserToken(user, conf.Cookies.SessionTimeout, conf.Jwt.SecretKey)
		if err != nil {
			return err
		}
		res.Token = tok
		return nil
	})
	return res.Token, err
}

// Auth Usecase를 반환합니다.
func NewAuth(userdb database.UserDatabase) AuthUsecase {
	return &AuthUC{userdb}
}
