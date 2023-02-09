package router

import (
	"github.com/JongGeonClass/JGC-API/database"
	"github.com/thak1411/gorn"
)

// 전체 EndPoint를 묶어서 제공합니다.
func New(
	userdb database.UserDatabase,
) *gorn.Router {
	router := gorn.NewRouter()

	auth := NewAuth(userdb)

	router.Extends("/api/auth", auth)

	return router
}
