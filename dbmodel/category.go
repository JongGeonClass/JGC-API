package dbmodel

import "github.com/thak1411/gorn"

// 미리 정의된 카테고리를 담아놓을 테이블입니다.
type Category struct {
	Id          int64  `rnsql:"id"  rntype:"INT"  rnopt:"PK NN UQ AI"  json:"id"`
	Name        string `rnsql:"name"  rntype:"VARCHAR(30)"  rnopt:"NN"  json:"name"`
	Description string `rnsql:"description"  rntype:"VARCHAR(200)"  rnopt:"NN"  json:"description"`
}

func init() {
	AddTable("CATEGORY", &Category{})
	AddIndex(&gorn.DBIndex{
		TableName: "CATEGORY",
		IndexName: "id_UNIQUE",
		IndexType: gorn.DBIndexTypeUnique,
		Columns: []*gorn.DBIndexColumn{
			{ColumnName: "id", ASC: true},
		},
	})
}
