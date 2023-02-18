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

// 개별 상품 정보 조회하기
func (h *ProductHandler) GetProduct(c *gorn.Context) {
	type Response struct { // 반환 타입
		Code    int                    `json:"code"`
		Product *dbmodel.PublicProduct `json:"product"`
	}
	ctx := c.GetContext()
	res := &Response{8000, nil}

	productId := c.GetParamInt64("product_id", -1) // 상품 번호를 가져옵니다.
	if err := c.Assert(productId >= 0, "id must be greater than or equal to 0"); err != nil {
		return
	}

	product, err := h.uc.GetProduct(ctx, productId) // 상품 정보를 가져옵니다.
	if err != nil {
		rnlog.Error("products get error: %+v", err)
		c.SendInternalServerError()
		return
	}
	res.Product = product
	c.SendJson(http.StatusOK, res)
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

// 장바구니에 담긴 상품 리스트를 가져옵니다.
// 이 함수는 항상 인증된 사용자만 사용할 수 있도록 미들웨어에서만 호출해야 합니다.
func (h *ProductHandler) GetCartProducts(c *gorn.Context) {
	type Response struct { // 반환 타입
		Code  int                   `json:"code"`
		Carts []*dbmodel.PublicCart `json:"carts"`
	}
	ctx := c.GetContext()
	res := &Response{8000, nil}
	conf := config.Get()
	token := c.GetValue(conf.Cookies.SessionName).(model.AuthUserTokenClaims)

	// 장바구니에 담긴 상품 리스트를 가져옵니다.
	carts, err := h.uc.GetCartProducts(ctx, token.Id)
	if err != nil {
		rnlog.Error("products get error: %+v", err)
		c.SendInternalServerError()
		return
	}
	res.Carts = carts
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

// 장바구니에 담긴 상품의 개수를 변경합니다.
// 이 함수는 항상 인증된 사용자만 사용할 수 있도록 미들웨어에서만 호출해야 합니다.
func (h *ProductHandler) UpdateCartAmount(c *gorn.Context) {
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
	// 장바구니에 상품 개수를 변경하는 로직을 실행합니다.
	if err := h.uc.UpdateCartAmount(ctx, token.Id, body.ProductId, body.Amount); err != nil {
		rnlog.Error("update cart amount error: %+v", err)
		c.SendInternalServerError()
		return
	}
	c.SendJson(http.StatusOK, res)
}

// 장바구니에 등록된 상품 삭제
// 이 함수는 항상 인증된 사용자만 사용할 수 있도록 미들웨어에서만 호출해야 합니다.
func (h *ProductHandler) DeleteFromCart(c *gorn.Context) {
	type Response struct { // 반환 타입
		Code int `json:"code"`
	}
	type Body struct { // Body 파라미터 타입
		ProductId int64 `json:"product_id"`
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
	// 장바구니에 상품을 삭제하는 로직을 실행합니다.
	if err := h.uc.DeleteFromCart(ctx, token.Id, body.ProductId); err != nil {
		rnlog.Error("delete from cart error: %+v", err)
		c.SendInternalServerError()
		return
	}
	c.SendJson(http.StatusOK, res)
}

// 상품에 리뷰 작성
// 이 함수는 항상 인증된 사용자만 사용할 수 있도록 미들웨어에서만 호출해야 합니다.
func (h *ProductHandler) AddReview(c *gorn.Context) {
	type Response struct { // 반환 타입
		Code     int   `json:"code"`
		ReviewId int64 `json:"review_id"`
	}
	type Body struct { // Body 파라미터 타입
		ProductId      int64  `json:"product_id"`
		Score          int64  `json:"score"`
		Content        string `json:"content"`
		ParentReviewId int64  `json:"parent_review_id"`
	}
	res := &Response{8000, 0}
	ctx := c.GetContext()
	body := &Body{}
	conf := config.Get()
	token := c.GetValue(conf.Cookies.SessionName).(model.AuthUserTokenClaims)
	if err := c.BindJsonBody(body); err != nil { // 바디 바인딩
		return
	}
	if err := c.AssertInt64Range(body.Score, 1, 5); err != nil {
		return
	}
	if err := c.Assert(body.ProductId > 0, "product_id must be greater than 0"); err != nil {
		return
	}
	if err := c.AssertStrLen(body.Content, 5, 1000); err != nil {
		return
	}
	// 리뷰를 작성하는 로직을 실행합니다.
	if reviewId, err := h.uc.AddReview(ctx, token.Id, body.ProductId, body.Score, body.ParentReviewId, &body.Content); err != nil {
		rnlog.Error("add review error: %+v", err)
		c.SendInternalServerError()
		return
	} else if reviewId == -1 { // 리뷰를 작성하려는 상품이 존재하지 않습니다.
		res.Code = 8001
	} else if reviewId == -2 { // 대댓글을 달려고 하는 부모 리뷰가 존재하지 않습니다.
		res.Code = 8002
	} else {
		res.ReviewId = reviewId
	}
	c.SendJson(http.StatusOK, res)
}

// 리뷰 조회하기
func (h *ProductHandler) GetReviews(c *gorn.Context) {
	type Response struct { // 반환 타입
		Code    int                     `json:"code"`
		Reviews []*dbmodel.PublicReview `json:"reviews"`
	}
	res := &Response{8000, nil}
	ctx := c.GetContext()
	productId := c.GetParamInt64("product_id", 0)
	if err := c.Assert(productId > 0, "product_id must be greater than 0"); err != nil {
		return
	}
	// 리뷰를 조회하는 로직을 실행합니다.
	if reviews, err := h.uc.GetReviews(ctx, productId); err != nil {
		rnlog.Error("get reviews error: %+v", err)
		c.SendInternalServerError()
		return
	} else {
		res.Reviews = reviews
	}
	c.SendJson(http.StatusOK, res)
}

// Product Handler를 반환합니다.
func NewProduct(uc usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{uc}
}
