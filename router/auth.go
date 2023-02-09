package router

import (
	"github.com/JongGeonClass/JGC-API/database"
	"github.com/JongGeonClass/JGC-API/handler"
	"github.com/JongGeonClass/JGC-API/usecase"
	"github.com/thak1411/gorn"
)

// Auth 관련 EndPoint를 묶어서 제공합니다.
func NewAuth(userdb database.UserDatabase) *gorn.Router {
	router := gorn.NewRouter()

	uc := usecase.NewAuth(userdb)
	hd := handler.NewAuth(uc)

	router.Post("/signup", hd.SignUp)
	router.Post("/login", hd.Login)
	router.Post("/logout", hd.Logout)
	return router
}
