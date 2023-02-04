package router

import (
	"github.com/thak1411/gorn"
)

// 전체 EndPoint를 묶어서 제공합니다.
func New() *gorn.Router {
	router := gorn.NewRouter()

	return router
}
