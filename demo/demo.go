package demo

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/JongGeonClass/JGC-API/database"
	"github.com/JongGeonClass/JGC-API/dbmodel"
	"github.com/JongGeonClass/JGC-API/util"
	"github.com/thak1411/rnlog"
)

// 디비에 데모 데이터를 생성합니다.
// 이때 없는 데이터만 추가로 생성합니다.
func Generate(
	userdb database.UserDatabase,
	productdb database.ProductDatabase,
) error {
	ctx := context.Background()

	// 유저 데이터 생성
	// id: morgan
	// password: morgan

	// id: andrewmjk1
	// password: andrewmjk1

	// id: effect
	// password: effect
	users := []*dbmodel.User{
		{
			Email:    "morgan@gmail.com",
			Nickname: "Morgan",
			Username: "morgan",
			Password: util.Encrypt256("morgan", "demo_salt"),
			Salt:     "demo_salt",
		},
		{
			Email:    "andrewmjk1@gmail.com",
			Nickname: "Andrewmjk1",
			Username: "andrewmjk1",
			Password: util.Encrypt256("andrewmjk1", "demo_salt"),
			Salt:     "demo_salt",
		},
		{
			Email:    "effect@gmail.com",
			Nickname: "Effect",
			Username: "effect",
			Password: util.Encrypt256("effect", "demo_salt"),
			Salt:     "demo_salt",
		},
	}
	rnlog.Info("Generating demo users...")
	minUserId := int64(0)
	for i, v := range users {
		if id, err := userdb.AddUser(ctx, v); err != nil {
			rnlog.Error("Error while adding user: %v", err)
			return err
		} else if i == 0 {
			minUserId = id
		}
	}

	// 카테고리 데이터 생성
	rnlog.Info("Generating demo categories...")
	minCategoryId := int64(0)
	for i := 1; i <= 9; i++ {
		id, err := productdb.AddCategory(ctx, &dbmodel.Category{
			Name:        "종건급 카테고리" + strconv.Itoa(i),
			Description: "종건급 카테고리" + strconv.Itoa(i) + "에 대한 설명입니다.",
		})
		if i == 1 {
			minCategoryId = id
		}
		if err != nil {
			rnlog.Error("Error while adding category: %v", err)
			return err
		}
	}

	// 브랜드 데이터 생성
	rnlog.Info("Generating demo brands...")
	minBrandId := int64(0)
	for i := 1; i <= 9; i++ {
		id, err := productdb.AddBrand(ctx, &dbmodel.Brand{
			UserId: int64(i%2) + minUserId,
			Name:   "종건급 브랜드" + strconv.Itoa(i),
			Email:  "thak1411@gmail.com",
		})
		if i == 1 {
			minBrandId = id
		}
		if err != nil {
			rnlog.Error("Error while adding brand: %v", err)
			return err
		}
	}

	// 데모 상품 추가
	rnlog.Info("Generating demo products...")
	minProductId := int64(0)
	for i := 1; i <= 50; i++ {
		id, err := productdb.AddProduct(ctx, &dbmodel.Product{
			BrandId:       int64(i%9 + int(minBrandId)),
			Name:          "종건급 상품" + strconv.Itoa(i),
			Price:         rand.Int63(),
			Amount:        rand.Int63(),
			TitleImageS3:  fmt.Sprintf("https://jgc-product.s3.ap-northeast-2.amazonaws.com/title/%d.txt", i),
			DescriptionS3: fmt.Sprintf("https://jgc-product.s3.ap-northeast-2.amazonaws.com/description/%d.txt", i),
		})
		if i == 1 {
			minProductId = id
		}
		if err != nil {
			rnlog.Error("Error while adding product: %v", err)
			return err
		}
	}

	// 데모 리뷰 추가
	rnlog.Info("Generating demo reviews...")
	for i := 1; i <= 50; i++ {
		reviewId := int64(0)
		sumScore := int64(0)
		for j := 1; j <= 3; j++ {
			sc := rand.Int63n(5) + 1
			review := &dbmodel.Review{
				ProductId: int64(i + int(minProductId) - 1),
				UserId:    int64((i+j)%3) + minUserId,
				Score:     sc,
				Content:   "종건급 상품 진자 개지립니다. 꼭 사용해보세요!",
			}
			sumScore += sc
			if j == 2 {
				review.ParentReviewId = reviewId
			}
			id, err := productdb.AddReview(ctx, review)
			reviewId = id
			if err != nil {
				rnlog.Error("Error while adding review: %v", err)
				return err
			}
		}
		// 데모 통계 페이지 추가
		statistics := &dbmodel.ProductStatistics{
			ProductId:      int64(i + int(minProductId) - 1),
			ReviewCount:    3,
			SumReviewScore: sumScore,
			SoldQuantity:   0,
		}
		err := productdb.AddProductStatistics(ctx, statistics)
		if err != nil {
			rnlog.Error("Error while adding product statistics: %v", err)
			return err
		}
	}

	// 데모 통계 페이지 추가
	// rnlog.Info("Generating demo statistics pages...")
	// for i := 1; i <= 50; i++ {
	// 	err := productdb.AddProductStatistics(ctx, &dbmodel.ProductStatistics{
	// 		ProductId:      int64(i + int(minProductId) - 1),
	// 		ReviewCount:    0,
	// 		SumReviewScore: 0,
	// 		SoldQuantity:   0,
	// 	})
	// 	if err != nil {
	// 		rnlog.Error("Error while adding product statistics: %v", err)
	// 		return err
	// 	}
	// }

	// 데모 상품 카테고리 추가
	rnlog.Info("Generating demo product categories...")
	for i := 1; i <= 50; i++ {
		for j := 1; j <= 2; j++ {
			err := productdb.AddProductCategory(ctx, &dbmodel.ProductCategoryMap{
				ProductId:  int64(i + int(minProductId) - 1),
				CategoryId: int64((i+j)%9 + int(minCategoryId)),
			})
			if err != nil {
				rnlog.Error("Error while adding product category: %v", err)
				return err
			}
		}
	}

	return nil
}

// 디비에 존재하는 데모 데이터를 삭제합니다.
func Remove(
	userdb database.UserDatabase,
	productdb database.ProductDatabase,
) error {
	ctx := context.Background()

	// 상품 카테고리 데이터 삭제
	rnlog.Info("Removing demo product categories...")
	if err := productdb.DeleteAllProductCategoryMap(ctx); err != nil {
		rnlog.Error("Error while deleting product category: %v", err)
		return err
	}

	// 상품 통계 데이터 삭제
	rnlog.Info("Removing demo product statistics...")
	if err := productdb.DeleteAllProductStatistics(ctx); err != nil {
		rnlog.Error("Error while deleting product statistics: %v", err)
		return err
	}

	// 모든 리뷰 삭제
	rnlog.Info("Removing demo reviews...")
	if err := productdb.DeleteAllReviews(ctx); err != nil {
		rnlog.Error("Error while deleting review: %v", err)
		return err
	}

	// 상품 데이터 삭제
	rnlog.Info("Removing demo products...")
	if err := productdb.DeleteAllProducts(ctx); err != nil {
		rnlog.Error("Error while deleting product: %v", err)
		return err
	}

	// 브랜드 데이터 삭제
	rnlog.Info("Removing demo brands...")
	if err := productdb.DeleteAllBrands(ctx); err != nil {
		rnlog.Error("Error while deleting brand: %v", err)
		return err
	}

	// 카테고리 데이터 삭제
	rnlog.Info("Removing demo categories...")
	if err := productdb.DeleteAllCategories(ctx); err != nil {
		rnlog.Error("Error while deleting category: %v", err)
		return err
	}

	// 유저 데이터 삭제
	rnlog.Info("Removing demo users...")
	if err := userdb.DeleteAllUsers(ctx); err != nil {
		rnlog.Error("Error while deleting user: %v", err)
		return err
	}
	return nil
}
