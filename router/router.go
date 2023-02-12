package router

import (
	"github.com/JongGeonClass/JGC-API/database"
	"github.com/thak1411/gorn"
)

// 전체 EndPoint를 묶어서 제공합니다.
func New(
	userdb database.UserDatabase,
	productdb database.ProductDatabase,
) *gorn.Router {
	router := gorn.NewRouter()

	auth := NewAuth(userdb)
	product := NewProduct(userdb, productdb)

	router.Extends("/api/auth", auth)
	router.Extends("/api/product", product)

	return router
}
