package util

import (
	"errors"
	"time"

	"github.com/JongGeonClass/JGC-API/dbmodel"
	"github.com/JongGeonClass/JGC-API/model"
	"github.com/dgrijalva/jwt-go/v4"
)

// JWT 관련 모듈입니다.
// config.JwtSecretKey를 기반으로 암호화 하고 복호화 합니다.

// dbmodel.User를 기준으로 JWT Token을 발급합니다.
func CreateUserToken(user *dbmodel.User, expire time.Duration, secretKey string) (string, error) {
	at := model.AuthUserTokenClaims{
		Id:          user.Id,
		Uuid:        NewUuid(),
		Nickname:    user.Nickname,
		CreatedTime: user.CreatedTime,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(expire)),
		},
	}
	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, &at)
	token, err := atoken.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

// dbmodel.User를 기준으로 JWT 토큰을 인증합니다.
func AuthUserToken(token string, secretKey string) (*jwt.Token, interface{}, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			ErrUnexpectedSigningMethod := errors.New("unexpected signing method")
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(secretKey), nil
	}
	claims := &model.AuthUserTokenClaims{}
	tok, err := jwt.ParseWithClaims(token, claims, keyFunc)
	return tok, *claims, err
}
