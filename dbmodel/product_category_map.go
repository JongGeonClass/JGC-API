package dbmodel

import "github.com/thak1411/gorn"

// 한 상품(프로덕트)에 종속되어있는 카테고리를 담은 N:M 맵입니다.
type ProductCategoryMap struct {
	ProductId  int64 `rnsql:"product_id"  rntype:"INT"  rnopt:"PK NN"  FK:"PRODUCT.id"  json:"product_id"`
	CategoryId int64 `rnsql:"category_id"  rntype:"INT"  rnopt:"PK NN"  FK:"CATEGORY.id"  json:"category_id"`
}

func init() {
	AddTable("PRODUCT_CATEGORY_MAP", &ProductCategoryMap{})
	AddIndex(&gorn.DBIndex{
		TableName: "PRODUCT_CATEGORY_MAP",
		IndexName: "id_UNIQUE",
		IndexType: gorn.DBIndexTypeUnique,
		Columns: []*gorn.DBIndexColumn{
			{ColumnName: "product_id", ASC: true},
			{ColumnName: "category_id", ASC: true},
		},
	})
}
