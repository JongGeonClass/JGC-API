package router

import (
	"github.com/JongGeonClass/JGC-API/database"
	"github.com/JongGeonClass/JGC-API/handler"
	"github.com/JongGeonClass/JGC-API/usecase"
	"github.com/thak1411/gorn"
)

// Product 관련 EndPoint를 묶어서 제공합니다.
func NewProduct(
	userdb database.UserDatabase,
	productdb database.ProductDatabase,
) *gorn.Router {
	router := gorn.NewRouter()

	uc := usecase.NewProduct(userdb, productdb)
	hd := handler.NewProduct(uc)

	router.Get("/products", hd.GetProducts)

	return router
}
