package dbmodel

import "github.com/thak1411/gorn"

// 상품에 대한 통계를 저장해주는 테이블입니다.
type ProductStatistics struct {
	ProductId      int64 `rnsql:"product_id"  rntype:"INT"  rnopt:"PK NN"  FK:"PRODUCT.id"  json:"product_id"`
	ReviewCount    int64 `rnsql:"review_count"  rntype:"INT"  rnopt:"NN"  json:"review_count"`
	SumReviewScore int64 `rnsql:"sum_review_score"  rntype:"INT"  rnopt:"NN"  json:"sum_review_score"`
	SoldQuantity   int64 `rnsql:"sold_quantity"  rntype:"INT"  rnopt:"NN"  json:"sold_quantity"`
}

func init() {
	AddTable("PRODUCT_STATISTICS", &ProductStatistics{})
	AddIndex(&gorn.DBIndex{
		TableName: "PRODUCT_STATISTICS",
		IndexName: "id_UNIQUE",
		IndexType: gorn.DBIndexTypeUnique,
		Columns: []*gorn.DBIndexColumn{
			{ColumnName: "product_id", ASC: true},
		},
	})
}
