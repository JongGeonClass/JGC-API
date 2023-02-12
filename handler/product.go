package handler

import (
	"net/http"

	"github.com/JongGeonClass/JGC-API/dbmodel"
	"github.com/JongGeonClass/JGC-API/usecase"
	"github.com/thak1411/gorn"
	"github.com/thak1411/rnlog"
)

// Product Hanlder의 구현체입니다.
type ProductHandler struct {
	uc usecase.ProductUsecase
}

// 상품 리스트 조회하기
func (h *ProductHandler) GetProducts(c *gorn.Context) {
	type Response struct { // 반환 타입
		Code     int                      `json:"code"`
		Products []*dbmodel.PublicProduct `json:"products"`
	}
	ctx := c.GetContext()
	res := &Response{8000, nil}

	page := c.GetParamInt("page", 0) // 검색할 페이지 번호를 가져옵니다.
	if err := c.Assert(page >= 0, "page must be greater than or equal to 0"); err != nil {
		return
	}
	pagesize := c.GetParamInt("pagesize", -1) // 검색할 페이지 길이를 가져옵니다.
	if err := c.AssertIntRange(pagesize, 1, 100); err != nil {
		return
	}

	// 상품 리스트를 가져옵니다.
	products, err := h.uc.GetProducts(ctx, page, pagesize)
	if err != nil {
		rnlog.Error("products get error: %+v", err)
		c.SendInternalServerError()
		return
	}
	res.Products = products
	c.SendJson(http.StatusOK, res)
}

// Product Handler를 반환합니다.
func NewProduct(uc usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{uc}
}
