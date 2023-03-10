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

	router.Get("/product", hd.GetProduct)
	router.Get("/products", hd.GetProducts)
	router.Get("/carts", decode, hd.GetCartProducts)
	router.Post("/add-to-cart", decode, hd.AddToCart)
	router.Post("/update-cart-amount", decode, hd.UpdateCartAmount)
	router.Delete("/delete-cart-product", decode, hd.DeleteFromCart)
	router.Post("/add-review", decode, hd.AddReview)
	router.Get("/reviews", hd.GetReviews)
	router.Get("/categories", hd.GetCategories)
	router.Post("/add-pbv", decode, hd.AddPbvOption)
	router.Get("/pbv", decode, hd.GetPbvOption)
	router.Post("update-pbv", decode, hd.UpdatePbvOption)
	router.Delete("/delete-pbv", decode, hd.DeletePbvOption)
	router.Get("/brands", decode, hd.GetBrands)

	return router
}
