package demo

import (
	"context"
	"fmt"
	"math/rand"

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
			Id:       1,
			Email:    "root@mobis.com",
			Nickname: "Mobis",
			Username: "mobis",
			Password: util.Encrypt256("mobis", "demo_salt"),
			Salt:     "demo_salt",
		},
		{
			Id:       2,
			Email:    "root@hyundai.com",
			Nickname: "Hyundai",
			Username: "hyundai",
			Password: util.Encrypt256("hyundai", "demo_salt"),
			Salt:     "demo_salt",
		},
		{
			Id:       3,
			Email:    "root@kia.com",
			Nickname: "Kia",
			Username: "kia",
			Password: util.Encrypt256("kia", "demo_salt"),
			Salt:     "demo_salt",
		},
		{
			Id:       4,
			Email:    "morgan@gmail.com",
			Nickname: "Morgan",
			Username: "morgan",
			Password: util.Encrypt256("morgan", "demo_salt"),
			Salt:     "demo_salt",
		},
		{
			Id:       5,
			Email:    "andrewmjk1@gmail.com",
			Nickname: "Andrewmjk1",
			Username: "andrewmjk1",
			Password: util.Encrypt256("andrewmjk1", "demo_salt"),
			Salt:     "demo_salt",
		},
		{
			Id:       6,
			Email:    "effect@gmail.com",
			Nickname: "Effect",
			Username: "effect",
			Password: util.Encrypt256("effect", "demo_salt"),
			Salt:     "demo_salt",
		},
	}
	rnlog.Info("Generating demo users...")
	for _, v := range users {
		if _, err := userdb.AddUser(ctx, v); err != nil {
			rnlog.Error("Error while adding user: %v", err)
			return err
		}
	}

	// 카테고리 데이터 생성
	rnlog.Info("Generating demo categories...")
	categoryName := []string{"허드(HUD)", "네비게이션", "음향 기기", "의자", "카페트", "시트", "핸들", "방향제", "소품"}
	for i := 1; i <= 9; i++ {
		if _, err := productdb.AddCategory(ctx, &dbmodel.Category{
			Id:          int64(i),
			Name:        categoryName[i-1],
			Description: categoryName[i-1] + "에 대한 설명입니다.",
		}); err != nil {
			rnlog.Error("Error while adding category: %v", err)
			return err
		}
	}

	// 브랜드 데이터 생성
	rnlog.Info("Generating demo brands...")
	brandOwner := []int64{1, 2, 3, 1}
	brandName := []string{"현대 모비스(Hyundai Mobis)", "현대(Hyundai)", "기아(Kia)", "현대 몹이스(Hyundai Mopis)"}
	brandEmail := []string{
		"root@mobis.com", "root@hyundai.com", "root@kia.com", "root@mobis.com",
	}
	for i := 1; i <= 4; i++ {
		if _, err := productdb.AddBrand(ctx, &dbmodel.Brand{
			Id:     int64(i),
			UserId: brandOwner[i-1],
			Name:   brandName[i-1],
			Email:  brandEmail[i-1],
		}); err != nil {
			rnlog.Error("Error while adding brand: %v", err)
			return err
		}
	}

	// 데모 상품 추가
	rnlog.Info("Generating demo products...")
	brandId := []int64{
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 3, 3,
		3, 3, 3, 3,
		3, 3, 3, 3,
	}
	productName := []string{
		"Mobis 4K HUD - VK473824", "Mobis Super Fast Navigation - SUR39482", "Mobis High Quality Speaker - FJJS48374", "Mobis Miller Chair - FN2847184",
		"Mobis Very Soft Carpet - SHJF284719", "Mobis Ultra Comfortable Sheat - AOT38972", "Mobis Best Driver Hadle - VJG837592", "Mobis Malon Diffuser - VJG837592",
		"Mobis Cool Air Freshener - VJG837592", "Mobis Logo Sticker - KLIFJ3847294", "Kia 4K HUD - FJ24783", "Kia Super Fast Navigation - FJSF384571",
		"Kia High Quality Speaker - JFIE385729", "Kia Miller Chair - JGKS38472", "Kia Very Soft Carpet - JHIKS348724", "Kia Ultra Comfortable Sheat - GFJ3548724",
		"Kia Best Driver Hadle - SFD2587", "Kia Malon Diffuser - JGIS58276", "Kia Cool Air Freshener - JHI34872", "Kia Logo Sticker - OITU39571",
	}
	for i := 1; i <= 20; i++ {
		if _, err := productdb.AddProduct(ctx, &dbmodel.Product{
			Id:            int64(i),
			BrandId:       brandId[i-1],
			Name:          productName[i-1],
			Price:         int64(rand.Int31n(1000000) + 10000),
			Amount:        int64(rand.Int31n(10000) + 1),
			TitleImageS3:  fmt.Sprintf("https://jgc-product-bucket.s3.us-east-2.amazonaws.com/title/%d.png", i),
			DescriptionS3: fmt.Sprintf("https://jgc-product-bucket.s3.us-east-2.amazonaws.com/description/%d.txt", i),
		}); err != nil {
			rnlog.Error("Error while adding product: %v", err)
			return err
		}
	}

	// 데모 리뷰 추가
	rnlog.Info("Generating demo reviews...")
	reviewUser := []int64{1, 2, 3}
	reviewContent := []string{"좋은 상품 배달 잘 받았습니다.\n포장 상태도 양호하고 배달도 아주 빠르게 잘 도착했습니다.\n\n제품 퀄리티도 매우 좋아서 잘 사용하겠습니다. 감사합니다.",
		"감사합니다. 잘 사용하겠습니다.", "test"}
	for i := 1; i <= 20; i++ {
		reviewId := int64(0)
		sumScore := int64(0)
		for j := 1; j <= 3; j++ {
			sc := rand.Int63n(5) + 3
			if j == 1 {
				sc = 5
			}
			review := &dbmodel.Review{
				Id:        int64((i-1)*3 + j),
				ProductId: int64(i),
				UserId:    reviewUser[j-1],
				Score:     sc,
				Content:   reviewContent[j-1],
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
			ProductId:      int64(i),
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
	// for i := 1; i <= 20; i++ {
	// 	err := productdb.AddProductStatistics(ctx, &dbmodel.ProductStatistics{
	// 		ProductId:      int64(i),
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
	categoryId := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9}
	for i := 1; i <= 20; i++ {
		err := productdb.AddProductCategory(ctx, &dbmodel.ProductCategoryMap{
			ProductId:  int64(i),
			CategoryId: categoryId[i-1],
		})
		if err != nil {
			rnlog.Error("Error while adding product category: %v", err)
			return err
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
