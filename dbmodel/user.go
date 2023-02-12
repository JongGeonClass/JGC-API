package dbmodel

import (
	"time"

	"github.com/thak1411/gorn"
)

// 유저 정보를 담은 테이블입니다.
type User struct {
	Id          int64     `rnsql:"id"  rntype:"INT"  rnopt:"PK NN UQ AI"  json:"id"`
	Email       string    `rnsql:"email"  rntype:"VARCHAR(200)"  rnopt:"NN"  json:"email"`
	Nickname    string    `rnsql:"nickname"  rntype:"VARCHAR(30)"  rnopt:"NN"  json:"nickname"`
	Username    string    `rnsql:"username"  rntype:"VARCHAR(30)"  rnopt:"NN"  json:"username"`
	Password    string    `rnsql:"password"  rntype:"VARCHAR(512)"  rnopt:"NN"  json:"password"`
	Salt        string    `rnsql:"salt"  rntype:"VARCHAR(512)"  rnopt:"NN"  json:"salt"`
	CreatedTime time.Time `rnsql:"created_time"  rntype:"DATETIME"  rnopt:"NN"  json:"created_time"`
	UpdatedTime time.Time `rnsql:"updated_time"  rntype:"DATETIME"  rnopt:"NN"  json:"updated_time"`
}

func init() {
	AddTable("USER", &User{})
	AddIndex(&gorn.DBIndex{
		TableName: "USER",
		IndexName: "id_UNIQUE",
		IndexType: gorn.DBIndexTypeUnique,
		Columns: []*gorn.DBIndexColumn{
			{ColumnName: "id", ASC: true},
		},
	})
	AddIndex(&gorn.DBIndex{
		TableName: "USER",
		IndexName: "username_UNIQUE",
		IndexType: gorn.DBIndexTypeUnique,
		Columns: []*gorn.DBIndexColumn{
			{ColumnName: "username", ASC: true},
		},
	})
	AddIndex(&gorn.DBIndex{
		TableName: "USER",
		IndexName: "nickname_INDEX",
		IndexType: gorn.DBIndexTypeIndex,
		Columns: []*gorn.DBIndexColumn{
			{ColumnName: "nickname", ASC: true},
		},
	})
}
