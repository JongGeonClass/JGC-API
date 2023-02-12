package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/JongGeonClass/JGC-API/config"
	"github.com/JongGeonClass/JGC-API/usecase"
	"github.com/thak1411/gorn"
	"github.com/thak1411/rnlog"
)

// Auth Hanlder의 구현체입니다.
type AuthHandler struct {
	uc usecase.AuthUsecase
}

// 계정 회원가입
func (h *AuthHandler) SignUp(c *gorn.Context) {
	type Response struct { // 반환 타입
		Code int `json:"code"`
	}
	type Body struct { // Body 파라미터 타입
		Email    string `json:"email"`
		Nickname string `json:"nickname"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	reg := []string{
		"^.+@.+\\..+$",             // Email Regex
		"^[a-zA-Z0-9ㄱ-ㅎ가-힣-ㅏ-ㅣ]+$", // Nickname Regex
		"^[a-zA-Z0-9]+$",           // Username Regex
		"^[a-zA-Z0-9\\~\\!\\@\\#\\$\\%\\^\\&\\*\\(\\)\\-\\_\\+\\=\\[\\]\\{\\}\\.\\,\\<\\>\\/\\?\\;\\:\\'\\\"\\\\\\|\\`]+$", // Password Regex
	}
	lenLR := [][]int{
		{1, 100}, // Email Length
		{4, 30},  // Nickname Length
		{6, 30},  // Username Length
		{6, 30},  // Password Length
	}
	// 만약 비밀번호 길이 등 입력 제한 사항이 바뀌었다면 로그인 시 정규식 체크에서 혼란이 생기지 않도록 조심해야 합니다.
	res := &Response{8000}
	body := &Body{}
	if err := c.BindJsonBody(body); err != nil { // 바디 바인딩
		return
	}
	bodyItems := []string{body.Email, body.Nickname, body.Username, body.Password}
	for i, v := range bodyItems {
		if err := c.AssertStrLen(v, lenLR[i][0], lenLR[i][1]); err != nil {
			return
		}
	}
	// Check Regex for Email, Nickname, Username, Password
	for i, v := range bodyItems {
		if err := c.AssertStrRegex(v, reg[i]); err != nil {
			return
		}
	}

	// 회원 가입 로직을 실행합니다.
	ctx := c.GetContext()
	userId, err := h.uc.SignUp(ctx, body.Email, body.Nickname, body.Username, body.Password)
	if err != nil {
		rnlog.Error("SignUp error: %+v", err)
		c.SendInternalServerError()
		return
	} else if userId == -1 {
		res.Code = 8001
	} else if userId == -2 {
		res.Code = 8002
	}
	c.SendJson(http.StatusOK, res)
}

// 로그인을 시도합니다.
// 성공하면 토큰을 발급하고, 실패하면 에러를 반환합니다.
// 로그인에 성공했을 시에 쿠키를 생성합니다.
// 토큰과 쿠키의 유효기간은 일주일로 설정합니다.
func (h *AuthHandler) Login(c *gorn.Context) {
	type Response struct { // 반환 타입
		Code  int    `json:"code"`
		Token string `json:"token"`
	}
	type Body struct { // Body 파라미터 타입
		Username string `json:"username"`
		Password string `json:"password"`
	}
	//TODO: 정규식을 제대로 작성해야 합니다.
	reg := []string{
		"^.*$", // Username Regex
		"^.*$", // Password Regex
		// 비밀번호 정규식이 바뀌었다면 ChangePassword에서도 바꿔야 합니다.
	}
	res := &Response{8000, ""}
	conf := config.Get()
	body := &Body{}
	if err := c.BindJsonBody(body); err != nil {
		return
	}
	// 정규식 체크를 합니다. 만약 비밀번호 길이 등 입력 제한 사항이 바뀌었더라도 통과될 수 있어야 합니다.
	// 따라서 길이 검사만 하는 것을 권장합니다.
	for i, v := range []string{body.Username, body.Password} {
		if err := c.AssertStrRegex(v, reg[i]); err != nil {
			return
		}
	}

	// 로그인 로직을 실행합니다.
	ctx := c.GetContext()
	token, err := h.uc.Login(ctx, body.Username, body.Password)
	if err != nil {
		rnlog.Error("Login error: %+v", err)
		c.SendInternalServerError()
		return
	} else if token != "" {
		// 쿠키를 설정합니다.
		cookie := &http.Cookie{
			Name:     conf.Cookies.SessionName,
			Domain:   conf.Domain,
			Path:     "/",
			Expires:  time.Now().Add(conf.Cookies.SessionTimeout), // 쿠키 유효기간은 일주일입니다.
			Value:    token,
			HttpOnly: true,
		}
		// 퍼블릿 유저 쿠키를 설정합니다.
		// 이 쿠키는 유저의 정보를 가져올 때 사용합니다.
		cookie2 := &http.Cookie{
			Name:    conf.Cookies.PublicSessionName,
			Domain:  conf.Domain,
			Path:    "/",
			Expires: time.Now().Add(conf.Cookies.SessionTimeout), // 쿠키 유효기간은 일주일입니다.
			Value:   strings.Split(token, ".")[1],
		}
		c.SetCookie(cookie)
		c.SetCookie(cookie2)
		res.Token = token
	} else {
		res.Code = 8001
	}
	c.SendJson(http.StatusOK, res)
}

// 로그아웃 시킵니다.
// 브라우저에 설정된 쿠키를 삭제합니다.
func (h *AuthHandler) Logout(c *gorn.Context) {
	type Response struct { // 반환 타입
		Code int `json:"code"`
	}
	res := &Response{8000}
	conf := config.Get()

	// 쿠키를 삭제합니다.
	cookie := &http.Cookie{
		Name:   conf.Cookies.SessionName,
		Domain: conf.Domain,
		Path:   "/",
		MaxAge: -1,
	}
	// 퍼블릭 쿠키도 삭제합니다.
	cookie2 := &http.Cookie{
		Name:   conf.Cookies.PublicSessionName,
		Domain: conf.Domain,
		Path:   "/",
		MaxAge: -1,
	}
	c.SetCookie(cookie)
	c.SetCookie(cookie2)
	c.SendJson(http.StatusOK, res)
}

// Auth Handler를 반환합니다.
func NewAuth(uc usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc}
}
