package demo

import (
	"context"
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
	for _, v := range users {
		if _, err := userdb.AddUser(ctx, v); err != nil {
			rnlog.Error("Error while adding user: %v", err)
			return err
		}
	}

	// 브랜드 데이터 생성
	rnlog.Info("Generating demo brands...")
	minBrandId := int64(0)
	for i := 1; i <= 9; i++ {
		id, err := productdb.AddBrand(ctx, &dbmodel.Brand{
			Name:  "브랜드" + strconv.Itoa(i),
			Email: "thak1411@gmail.com",
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
	for i := 1; i <= 50; i++ {
		_, err := productdb.AddProduct(ctx, &dbmodel.Product{
			BrandId:       int64(i%9 + int(minBrandId)),
			Name:          "데모 상품" + strconv.Itoa(i),
			Price:         rand.Int63(),
			Amount:        rand.Int63(),
			TitleImageS3:  "demo_image_s3_link",
			DescriptionS3: "demo_description_s3_link",
		})
		if err != nil {
			rnlog.Error("Error while adding product: %v", err)
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

	// 유저 데이터 삭제
	rnlog.Info("Removing demo users...")
	if err := userdb.DeleteAllUsers(ctx); err != nil {
		rnlog.Error("Error while deleting user: %v", err)
		return err
	}
	return nil
}