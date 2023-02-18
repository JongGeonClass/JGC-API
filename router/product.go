package router

import (
	"github.com/JongGeonClass/JGC-API/database"
	"github.com/JongGeonClass/JGC-API/handler"
	"github.com/JongGeonClass/JGC-API/middleware"
	"github.com/JongGeonClass/JGC-API/usecase"
	"github.com/thak1411/gorn"
)

// Product 관련 EndPoint를 묶어서 제공합니다.
func NewProduct(
	userdb database.UserDatabase,
	productdb database.ProductDatabase,
) *gorn.Router {
	router := gorn.NewRouter()

	md := middleware.NewAuth(userdb)
	uc := usecase.NewProduct(userdb, productdb)
	hd := handler.NewProduct(uc)

	decode := md.TokenDecode

	router.Get("/products", hd.GetProducts)
	router.Post("/add-to-cart", decode, hd.AddToCart)
	router.Post("/update-cart-amount", decode, hd.UpdateCartAmount)
	router.Delete("/delete-cart-product", decode, hd.DeleteFromCart)
	router.Post("/add-review", decode, hd.AddReview)

	return router
}
