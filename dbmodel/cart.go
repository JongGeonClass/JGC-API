package dbmodel

import (
	"time"

	"github.com/thak1411/gorn"
)

// 유저가 상품(프로덕트)을 담아놓은 장바구니 정보를 담은 N:M 맵입니다.
type Cart struct {
	UserId      int64     `rnsql:"user_id"  rntype:"INT"  rnopt:"NN"  FK:"USER.id"  json:"user_id"`
	ProductId   int64     `rnsql:"product_id"  rntype:"INT"  rnopt:"NN"  FK:"PRODUCT.id"  json:"product_id"`
	Amount      int64     `rnsql:"amount"  rntype:"INT"  rnopt:"NN"  json:"amount"`
	CreatedTime time.Time `rnsql:"created_time"  rntype:"DATETIME"  rnopt:"NN"  json:"created_time"`
	UpdatedTime time.Time `rnsql:"updated_time"  rntype:"DATETIME"  rnopt:"NN"  json:"updated_time"`
}

func init() {
	AddTable("CART", &Cart{})
	AddIndex(&gorn.DBIndex{
		TableName: "CART",
		IndexName: "id_UNIQUE",
		IndexType: gorn.DBIndexTypeUnique,
		Columns: []*gorn.DBIndexColumn{
			{ColumnName: "user_id", ASC: true},
			{ColumnName: "product_id", ASC: true},
		},
	})
}
