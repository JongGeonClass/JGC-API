package model

import (
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

// dbmodel.User를 기반으로 인증하는 토큰에 포함될 정보들입니다.
type AuthUserTokenClaims struct {
	Id          int64     `json:"id"`
	Uuid        string    `json:"uuid"`
	Nickname    string    `json:"nickname"`
	CreatedTime time.Time `json:"created_time"`
	jwt.StandardClaims
}

// 게스트 토큰입니다.
var GuestAuthUserTokenClaims = AuthUserTokenClaims{
	Id:          -2,
	Uuid:        "",
	Nickname:    "",
	CreatedTime: time.Time{},
}
