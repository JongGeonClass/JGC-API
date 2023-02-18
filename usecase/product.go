package usecase

import (
	"context"

	"github.com/JongGeonClass/JGC-API/database"
	"github.com/JongGeonClass/JGC-API/dbmodel"
)

// Product Usecase의 인터페이스입니다.
type ProductUsecase interface {
	GetProducts(ctx context.Context, page, pagesize int) ([]*dbmodel.PublicProduct, int, error)
	AddToCart(ctx context.Context, userId, productId, amount int64) error
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

// 장바구니에 상품을 추가합니다.
// 존재하지 않는 상품이라면 무시합니다.
// 장바구니에 이미 상품이 담겨있다면, 기존의 개수에 추가로 개수를 더해줍니다.
func (uc *ProductUC) AddToCart(ctx context.Context, userId, productId, amount int64) error {
	// 존재하는 상품인지 확인합니다.
	if exists, err := uc.productdb.CheckProductExists(ctx, productId); err != nil {
		return err
	} else if !exists {
		// 존재하지 않는 상품이라면 무시합니다.
		return nil
	}

	err := uc.productdb.ExecTx(ctx, func(txdb database.ProductDatabase) error {
		// 장바구니에 이미 상품이 담겨있는지 확인합니다.
		if isExists, err := txdb.CheckCartHasProduct(ctx, userId, productId); err != nil {
			return err
		} else if !isExists {
			// 존재하지 않는다면 추가하고 종료합니다.
			if err := txdb.AddCart(ctx, &dbmodel.Cart{
				UserId:    userId,
				ProductId: productId,
				Amount:    amount,
			}); err != nil {
				return err
			}
			return nil
		}
		// 존재한다면, 개수를 더해줍니다.
		cart, err := txdb.GetCartProduct(ctx, userId, productId)
		if err != nil {
			return err
		}
		cart.Amount += amount
		// 업데이트 된 개수를 반영합니다.
		return txdb.UpdateCart(ctx, cart)
	})
	return err
}

// 장바구니에 상품을 삭제합니다.

// Product Usecase를 반환합니다.
func NewProduct(
	userdb database.UserDatabase,
	productdb database.ProductDatabase,
) ProductUsecase {
	return &ProductUC{userdb, productdb}
}
