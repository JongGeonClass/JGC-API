package dbmodel

import (
	"time"

	"github.com/thak1411/gorn"
)

// 유저가 작성한 리뷰를 담은 테이블입니다.
type Review struct {
	Id             int64     `rnsql:"id"  rntype:"INT"  rnopt:"PK NN UQ AI"  json:"id"`
	ProductId      int64     `rnsql:"product_id"  rntype:"INT"  rnopt:"NN"  FK:"PRODUCT.id"  json:"product_id"`
	UserId         int64     `rnsql:"user_id"  rntype:"INT"  rnopt:"NN"  FK:"USER.id"  json:"user_id"`
	Score          int64     `rnsql:"score"  rntype:"INT"  rnopt:"NN"  json:"score"`
	Content        string    `rnsql:"content"  rntype:"VARCHAR(1000)"  rnopt:"NN"  json:"content"`
	ParentReviewId int64     `rnsql:"parent_review_id"  rntype:"INT"  rnopt:"NN"  json:"parent_review_id"`
	CreatedTime    time.Time `rnsql:"created_time"  rntype:"DATETIME"  rnopt:"NN"  json:"created_time"`
	UpdatedTime    time.Time `rnsql:"updated_time"  rntype:"DATETIME"  rnopt:"NN"  json:"updated_time"`
}

func init() {
	AddTable("REVIEW", &Review{})
	AddIndex(&gorn.DBIndex{
		TableName: "REVIEW",
		IndexName: "id_UNIQUE",
		IndexType: gorn.DBIndexTypeUnique,
		Columns: []*gorn.DBIndexColumn{
			{ColumnName: "id", ASC: true},
		},
	})
}
