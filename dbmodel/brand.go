package dbmodel

import (
	"time"

	"github.com/thak1411/gorn"
)

// 브랜드 정보를 담은 테이블입니다.
type Brand struct {
	Id          int64     `rnsql:"id"  rntype:"INT"  rnopt:"PK NN UQ AI"  json:"id"`
	Name        string    `rnsql:"name"  rntype:"VARCHAR(200)"  rnopt:"NN"  json:"name"`
	Email       string    `rnsql:"email"  rntype:"VARCHAR(200)"  rnopt:"NN"  json:"email"`
	CreatedTime time.Time `rnsql:"created_time"  rntype:"DATETIME"  rnopt:"NN"  json:"created_time"`
	UpdatedTime time.Time `rnsql:"updated_time"  rntype:"DATETIME"  rnopt:"NN"  json:"updated_time"`
}

func init() {
	AddTable("BRAND", &Brand{})
	AddIndex(&gorn.DBIndex{
		TableName: "BRAND",
		IndexName: "id_UNIQUE",
		IndexType: gorn.DBIndexTypeUnique,
		Columns: []*gorn.DBIndexColumn{
			{ColumnName: "id", ASC: true},
		},
	})
}
