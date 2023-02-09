package middleware

import (
	"github.com/JongGeonClass/JGC-API/config"
	"github.com/JongGeonClass/JGC-API/database"
	"github.com/JongGeonClass/JGC-API/model"
	"github.com/JongGeonClass/JGC-API/util"
	"github.com/thak1411/gorn"
)

type AuthMiddleware struct {
	userdb database.UserDatabase
}

// 토큰을 디코딩 하고 Request Context에 디코딩된 토큰 정보를 탑재합니다.
func (md *AuthMiddleware) TokenDecode(c *gorn.Context) {
	conf := config.Get()

	token, err := c.GetCookie(conf.Cookies.SessionName)
	if err != nil {
		c.SendNotAuthorized()
		return
	}
	tok, claims, err := util.AuthUserToken(token.Value, conf.Jwt.SecretKey)
	if err != nil || !tok.Valid {
		c.SendNotAuthorized()
		return
	}
	c.SetValue(conf.Cookies.SessionName, claims)
}

// TokenDecode와 작동방식이 동일하지만, 토큰을 인증하지 못했을 경우에 다릅니다.
// 인증할 토큰이 없을 경우 게스트 토큰을 탑재합니다.
// 토큰을 디코딩 하고 Request Context에 디코딩된 토큰 정보를 탑재합니다.
func (md *AuthMiddleware) TokenDecodeWithGuest(c *gorn.Context) {
	conf := config.Get()

	var claims interface{}

	token, err := c.GetCookie(conf.Cookies.SessionName)
	if err != nil {
		claims = model.GuestAuthUserTokenClaims
	} else {
		tok, clm, err := util.AuthUserToken(token.Value, conf.Jwt.SecretKey)
		if err != nil || !tok.Valid {
			claims = model.GuestAuthUserTokenClaims
		} else {
			claims = clm
		}
	}
	c.SetValue(conf.Cookies.SessionName, claims)
}

// Auth Middleware를 반환합니다.
func NewAuth(userdb database.UserDatabase) *AuthMiddleware {
	return &AuthMiddleware{
		userdb: userdb,
	}
}
