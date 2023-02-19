package usecase

import (
	"context"

	"github.com/JongGeonClass/JGC-API/database"
	"github.com/JongGeonClass/JGC-API/dbmodel"
)

// Product Usecase의 인터페이스입니다.
type ProductUsecase interface {
	GetProduct(ctx context.Context, productId int64) (*dbmodel.PublicProduct, error)
	GetProducts(ctx context.Context, page, pagesize int) ([]*dbmodel.PublicProduct, int, error)
	GetCartProducts(ctx context.Context, userId int64) ([]*dbmodel.PublicCart, error)
	AddToCart(ctx context.Context, userId, productId, amount int64) error
	UpdateCartAmount(ctx context.Context, userId, productId, amount int64) error
	DeleteFromCart(ctx context.Context, userId, productId int64) error
	AddReview(ctx context.Context, userId, productId, score, parentReviewId int64, content *string) (int64, error)
	GetReviews(ctx context.Context, productId int64) ([]*dbmodel.PublicReview, error)
	GetCategories(ctx context.Context) ([]*dbmodel.Category, error)
}

// Product Usecase의 구현체입니다.
type ProductUC struct {
	userdb    database.UserDatabase
	productdb database.ProductDatabase
}

// 개별 상품 정보를 가져옵니다.
func (uc *ProductUC) GetProduct(ctx context.Context, productId int64) (*dbmodel.PublicProduct, error) {
	return uc.productdb.GetPublicProduct(ctx, productId)
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

// 장바구니에 담긴 상품 리스트를 가져옵니다.
func (uc *ProductUC) GetCartProducts(ctx context.Context, userId int64) ([]*dbmodel.PublicCart, error) {
	return uc.productdb.GetCartProducts(ctx, userId)
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

// 장바구니에 담긴 상품의 개수를 변경합니다.
func (uc *ProductUC) UpdateCartAmount(ctx context.Context, userId, productId, amount int64) error {
	err := uc.productdb.ExecTx(ctx, func(txdb database.ProductDatabase) error {
		// 장바구니에 이미 상품이 담겨있는지 확인합니다.
		if isExists, err := txdb.CheckCartHasProduct(ctx, userId, productId); err != nil {
			return err
		} else if !isExists {
			// 존재하지 않는다면 무시합니다.
			return nil
		}
		// 존재한다면, 개수를 변경합니다..
		cart := &dbmodel.Cart{
			UserId:    userId,
			ProductId: productId,
			Amount:    amount,
		}
		// 업데이트 된 개수를 반영합니다.
		return txdb.UpdateCart(ctx, cart)
	})
	return err
}

// 장바구니에서 상품을 삭제합니다.
// 만약 장바구니에 상품이 없다면 무시합니다.
func (uc *ProductUC) DeleteFromCart(ctx context.Context, userId, productId int64) error {
	if exists, err := uc.productdb.CheckCartHasProduct(ctx, userId, productId); err != nil {
		return err
	} else if exists {
		return uc.productdb.DeleteCartProduct(ctx, userId, productId)
	}
	return nil
}

// 상품에 리뷰를 작성합니다.
// 이후 작성한 리뷰 아이디를 반환합니다.
// 존재하지 않는 상품이라면 -1을 반환합니다.
// 대댓글을 달려고 할 때 존재하지 않는 부모 리뷰라면, -2를 반환합니다.
func (uc *ProductUC) AddReview(ctx context.Context, userId, productId, score, parentReviewId int64, content *string) (int64, error) {
	res := int64(0)
	err := uc.productdb.ExecTx(ctx, func(txdb database.ProductDatabase) error {
		// 존재하는 상품인지 확인합니다.
		if exists, err := txdb.CheckProductExists(ctx, productId); err != nil {
			return err
		} else if !exists {
			// 존재하지 않는 상품이라면 -1을 반환합니다.
			res = -1
			return nil
		}
		// 부모 리뷰 아이디가 있다면 같은 상품 내에 존재하는 리뷰인지 확인합니다.
		if parentReviewId != 0 {
			if exists, err := txdb.CheckReviewExists(ctx, parentReviewId); err != nil {
				return err
			} else if !exists {
				// 존재하지 않는 리뷰라면 -2를 반환합니다.
				res = -2
				return nil
			}
		}

		// 리뷰를 등록합니다.
		if reviewId, err := txdb.AddReview(ctx, &dbmodel.Review{
			ProductId:      productId,
			UserId:         userId,
			Score:          score,
			Content:        *content,
			ParentReviewId: parentReviewId,
		}); err != nil {
			return err
		} else {
			res = reviewId
		}

		// 이후 통계 테이블을 업데이트합니다.
		statistics, err := txdb.GetProductStatistics(ctx, productId)
		if err != nil {
			return err
		}
		statistics.ReviewCount++
		statistics.SumReviewScore += score
		return txdb.UpdateProductStatistics(ctx, statistics)
	})
	return res, err
}

// 리뷰 리스트를 가져옵니다.
func (uc *ProductUC) GetReviews(ctx context.Context, productId int64) ([]*dbmodel.PublicReview, error) {
	return uc.productdb.GetReviewList(ctx, productId)
}

// 카테고리 리스트를 가져옵니다.
func (uc *ProductUC) GetCategories(ctx context.Context) ([]*dbmodel.Category, error) {
	return uc.productdb.GetAllCategories(ctx)
}

// Product Usecase를 반환합니다.
func NewProduct(
	userdb database.UserDatabase,
	productdb database.ProductDatabase,
) ProductUsecase {
	return &ProductUC{userdb, productdb}
}
