package database

import (
	"context"
	"time"

	"github.com/JongGeonClass/JGC-API/dbmodel"
	"github.com/thak1411/gorn"
	"github.com/thak1411/rnlog"
)

// 유저 디비의 인터페이스 입니다.
type UserDatabase interface {
	ExecTx(ctx context.Context, fn func(txdb UserDatabase) error) error
	AddUser(ctx context.Context, user *dbmodel.User) (int64, error)
	DeleteAllUsers(ctx context.Context) error
	UpdateUser(ctx context.Context, user *dbmodel.User) error
	CheckUserExistsByUsername(ctx context.Context, username string) (bool, error)
	CheckUserExistsByNickname(ctx context.Context, nickname string) (bool, error)
	CheckUserExistsById(ctx context.Context, userId int64) (bool, error)
	GetUserById(ctx context.Context, id int64) (*dbmodel.User, error)
	GetUserByUsername(ctx context.Context, username string) (*dbmodel.User, error)
	GetUserByNickname(ctx context.Context, nickname string) (*dbmodel.User, error)
}

// 유저 디비의 구현체입니다.
type UserDB struct {
	*gorn.DB
}

// 넘겨받은 함수로 트랜잭션을 실행합니다.
func (h *UserDB) ExecTx(ctx context.Context, fn func(txdb UserDatabase) error) error {
	txdb, err := h.DB.BeginTx(ctx)
	if err != nil {
		return err
	}
	newHandler := &UserDB{
		DB: txdb,
	}
	err = fn(newHandler)
	if err != nil {
		if rbErr := txdb.RollbackTx(); rbErr != nil {
			rnlog.Error("Rollback error: %v", rbErr)
			return rbErr
		}
		return err
	}
	return txdb.CommitTx()
}

// 유저를 유저 디비에 추가합니다.
// 이후 유저 id를 반환합니다.
func (h *UserDB) AddUser(ctx context.Context, user *dbmodel.User) (int64, error) {
	ntime := time.Now()
	user.CreatedTime = ntime
	user.UpdatedTime = ntime
	return h.InsertWithLastId(ctx, "USER", user)
}

// 모든 유저를 삭제합니다.
func (h *UserDB) DeleteAllUsers(ctx context.Context) error {
	sql := gorn.NewSql().
		DeleteFrom("USER").
		Where("id > ?", -1)
	res, err := h.Exec(ctx, sql)
	if err != nil {
		return err
	}
	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// 유저를 수정합니다.
func (h *UserDB) UpdateUser(ctx context.Context, user *dbmodel.User) error {
	user.UpdatedTime = time.Now()
	sql := gorn.NewSql().
		Update("USER", user).
		Where("id = ?", user.Id)
	result, err := h.Exec(ctx, sql)
	if err != nil {
		return err
	}
	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// 해당 username을 가진 유저가 있는지 체크합니다.
func (h *UserDB) CheckUserExistsByUsername(ctx context.Context, username string) (bool, error) {
	type UserCount struct {
		Count int64 `rnsql:"COUNT(*)"`
	}
	count := &UserCount{}
	sql := gorn.NewSql().
		Select(count).
		From("USER").
		Where("username = ?", username)
	row := h.QueryRow(ctx, sql)
	err := h.ScanRow(row, count)
	if err != nil {
		return false, err
	}
	return count.Count > 0, nil
}

// 해당 nickname을 가진 유저가 있는지 체크합니다.
func (h *UserDB) CheckUserExistsByNickname(ctx context.Context, nickname string) (bool, error) {
	type UserCount struct {
		Count int64 `rnsql:"COUNT(*)"`
	}
	count := &UserCount{}
	sql := gorn.NewSql().
		Select(count).
		From("USER").
		Where("nickname = ?", nickname)
	row := h.QueryRow(ctx, sql)
	err := h.ScanRow(row, count)
	if err != nil {
		return false, err
	}
	return count.Count > 0, nil
}

// 해당 id를 가진 유저가 있는지 체크합니다.
func (h *UserDB) CheckUserExistsById(ctx context.Context, userId int64) (bool, error) {
	type UserCount struct {
		Count int64 `rnsql:"COUNT(*)"`
	}
	count := &UserCount{}
	sql := gorn.NewSql().
		Select(count).
		From("USER").
		Where("id = ?", userId)
	row := h.QueryRow(ctx, sql)
	err := h.ScanRow(row, count)
	if err != nil {
		return false, err
	}
	return count.Count > 0, nil
}

// id로 유저의 정보를 가져옵니다.
func (h *UserDB) GetUserById(ctx context.Context, id int64) (*dbmodel.User, error) {
	result := &dbmodel.User{}
	sql := gorn.NewSql().
		Select(result).
		From("USER").
		Where("id = ?", id)

	row := h.QueryRow(ctx, sql)
	err := h.ScanRow(row, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// username으로 유저의 정보를 가져옵니다.
func (h *UserDB) GetUserByUsername(ctx context.Context, username string) (*dbmodel.User, error) {
	result := &dbmodel.User{}
	sql := gorn.NewSql().
		Select(result).
		From("USER").
		Where("username = ?", username)

	row := h.QueryRow(ctx, sql)
	err := h.ScanRow(row, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// nickname으로 유저의 정보를 가져옵니다.
func (h *UserDB) GetUserByNickname(ctx context.Context, nickname string) (*dbmodel.User, error) {
	result := &dbmodel.User{}
	sql := gorn.NewSql().
		Select(result).
		From("USER").
		Where("nickname = ?", nickname)

	row := h.QueryRow(ctx, sql)
	err := h.ScanRow(row, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 새로운 디비 객체를 연결합니다.
func NewUser(db *gorn.DB) UserDatabase {
	return &UserDB{
		DB: db,
	}
}
