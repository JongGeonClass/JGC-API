package handler

import (
	"net/http"

	"github.com/JongGeonClass/JGC-API/config"
	"github.com/JongGeonClass/JGC-API/dbmodel"
	"github.com/JongGeonClass/JGC-API/model"
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
		Code        int                      `json:"code"`
		Products    []*dbmodel.PublicProduct `json:"products"`
		MaxPagesize int                      `json:"max_pagesize"`
	}
	ctx := c.GetContext()
	res := &Response{8000, nil, 1}

	page := c.GetParamInt("page", 0) // 검색할 페이지 번호를 가져옵니다.
	if err := c.Assert(page >= 0, "page must be greater than or equal to 0"); err != nil {
		return
	}
	pagesize := c.GetParamInt("pagesize", -1) // 검색할 페이지 길이를 가져옵니다.
	if err := c.AssertIntRange(pagesize, 1, 100); err != nil {
		return
	}

	// 상품 리스트를 가져옵니다.
	products, maxPagesize, err := h.uc.GetProducts(ctx, page, pagesize)
	if err != nil {
		rnlog.Error("products get error: %+v", err)
		c.SendInternalServerError()
		return
	}
	res.Products = products
	res.MaxPagesize = maxPagesize
	c.SendJson(http.StatusOK, res)
}

// 장바구니에 상품을 담습니다.
// 이 함수는 항상 인증된 사용자만 사용할 수 있도록 미들웨어에서만 호출해야 합니다.
func (h *ProductHandler) AddToCart(c *gorn.Context) {
	type Response struct { // 반환 타입
		Code int `json:"code"`
	}
	type Body struct { // Body 파라미터 타입
		ProductId int64 `json:"product_id"`
		Amount    int64 `json:"amount"`
	}
	res := &Response{8000}
	ctx := c.GetContext()
	body := &Body{}
	conf := config.Get()
	token := c.GetValue(conf.Cookies.SessionName).(model.AuthUserTokenClaims)
	if err := c.BindJsonBody(body); err != nil { // 바디 바인딩
		return
	}
	if err := c.Assert(body.ProductId > 0, "product_id must be greater than 0"); err != nil {
		return
	}
	if err := c.Assert(body.Amount > 0, "amount must be greater than 0"); err != nil {
		return
	}

	// 장바구니에 상품 담는 로직을 실행합니다.
	if err := h.uc.AddToCart(ctx, token.Id, body.ProductId, body.Amount); err != nil {
		rnlog.Error("add to cart error: %+v", err)
		c.SendInternalServerError()
		return
	}
	c.SendJson(http.StatusOK, res)
}

// Product Handler를 반환합니다.
func NewProduct(uc usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{uc}
}
