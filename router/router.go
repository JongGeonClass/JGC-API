package router

import (
	"github.com/JongGeonClass/JGC-API/config"
	"github.com/JongGeonClass/JGC-API/database"
	"github.com/thak1411/gorn"
)

// 전체 EndPoint를 묶어서 제공합니다.
func New(
	userdb database.UserDatabase,
	productdb database.ProductDatabase,
) *gorn.Router {
	conf := config.Get()
	router := gorn.NewRouter()

	auth := NewAuth(userdb)
	product := NewProduct(userdb, productdb)

	router.Extends("/api/auth", auth)
	router.Extends("/api/product", product)

	options := &gorn.RouterOptions{
		AllowedOrigins:   conf.CorsOrigin,
		AllowedHeaders:   []string{"*"},
		MaxAge:           conf.MaxAge,
		AllowCredentials: true,
	}
	router.SetOptions(options)

	return router
}
