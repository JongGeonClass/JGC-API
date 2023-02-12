package usecase

import (
	"context"

	"github.com/JongGeonClass/JGC-API/database"
	"github.com/JongGeonClass/JGC-API/dbmodel"
)

// Product Usecase의 인터페이스입니다.
type ProductUsecase interface {
	GetProducts(ctx context.Context, page, pagesize int) ([]*dbmodel.PublicProduct, int, error)
}

// Product Usecase의 구현체입니다.
type ProductUC struct {
	userdb    database.UserDatabase
	productdb database.ProductDatabase
}

// 상품 리스트를 가져옵니다.
func (uc *ProductUC) GetProducts(ctx context.Context, page, pagesize int) ([]*dbmodel.PublicProduct, int, error) {
	products, err := uc.productdb.GetProducts(ctx, page, pagesize)
	if err != nil {
		return nil, 0, err
	}
	productsCount, err := uc.productdb.GetProductsCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	max := func(i, j int) int {
		if i > j {
			return i
		}
		return j
	}
	return products, max(productsCount-1, 0) / pagesize, nil
}

// Product Usecase를 반환합니다.
func NewProduct(
	userdb database.UserDatabase,
	productdb database.ProductDatabase,
) ProductUsecase {
	return &ProductUC{userdb, productdb}
}
