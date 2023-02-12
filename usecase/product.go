package usecase

import (
	"context"

	"github.com/JongGeonClass/JGC-API/database"
	"github.com/JongGeonClass/JGC-API/dbmodel"
)

// Product Usecase의 인터페이스입니다.
type ProductUsecase interface {
	GetProducts(ctx context.Context, page, pagesize int) ([]*dbmodel.PublicProduct, error)
}

// Product Usecase의 구현체입니다.
type ProductUC struct {
	userdb    database.UserDatabase
	productdb database.ProductDatabase
}

// 상품 리스트를 가져옵니다.
func (uc *ProductUC) GetProducts(ctx context.Context, page, pagesize int) ([]*dbmodel.PublicProduct, error) {
	return uc.productdb.GetProducts(ctx, page, pagesize)
}

// Product Usecase를 반환합니다.
func NewProduct(
	userdb database.UserDatabase,
	productdb database.ProductDatabase,
) ProductUsecase {
	return &ProductUC{userdb, productdb}
}
