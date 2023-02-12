package dbmodel

import (
	"time"

	"github.com/thak1411/gorn"
)

// 판매자가 판매할 상품 정보를 담은 테이블입니다.
type Product struct {
	Id            int64     `rnsql:"id"  rntype:"INT"  rnopt:"PK NN UQ AI"  json:"id"`
	BrandId       int64     `rnsql:"brand_id"  rntype:"INT"  rnopt:"NN"  FK:"BRAND.id"  json:"brand_id"`
	Name          string    `rnsql:"name"  rntype:"VARCHAR(200)"  rnopt:"NN"  json:"name"`
	Price         int64     `rnsql:"price"  rntype:"INT"  rnopt:"NN"  json:"price"`
	Amount        int64     `rnsql:"amount"  rntype:"INT"  rnopt:"NN"  json:"amount"`
	DescriptionS3 string    `rnsql:"description_s3"  rntype:"VARCHAR(200)"  rnopt:"NN"  json:"description_s3"`
	CreatedTime   time.Time `rnsql:"created_time"  rntype:"DATETIME"  rnopt:"NN"  json:"created_time"`
	UpdatedTime   time.Time `rnsql:"updated_time"  rntype:"DATETIME"  rnopt:"NN"  json:"updated_time"`
}

func init() {
	AddTable("PRODUCT", &Product{})
	AddIndex(&gorn.DBIndex{
		TableName: "PRODUCT",
		IndexName: "id_UNIQUE",
		IndexType: gorn.DBIndexTypeUnique,
		Columns: []*gorn.DBIndexColumn{
			{ColumnName: "id", ASC: true},
		},
	})
}
