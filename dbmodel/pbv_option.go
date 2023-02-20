package dbmodel

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/thak1411/gorn"
)

// 유저가 담아놓은 커스텀 PBV 옵션을 담은 테이블입니다.
type PbvOption struct {
	Id          int64     `rnsql:"id"  rntype:"INT"  rnopt:"PK NN UQ AI"  json:"id"`
	UserId      int64     `rnsql:"user_id"  rntype:"INT"  rnopt:"NN"  FK:"USER.id"  json:"user_id"`
	Data        DataJson  `rnsql:"data"  rntype:"JSON"  rnopt:"NN"  json:"data"  db:"data"`
	CreatedTime time.Time `rnsql:"created_time"  rntype:"DATETIME"  rnopt:"NN"  json:"created_time"`
	UpdatedTime time.Time `rnsql:"updated_time"  rntype:"DATETIME"  rnopt:"NN"  json:"updated_time"`
}

type DataJson map[string]interface{}

func (d *DataJson) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *DataJson) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), d)
}

func init() {
	AddTable("PBV_OPTION", &PbvOption{})
	AddIndex(&gorn.DBIndex{
		TableName: "PBV_OPTION",
		IndexName: "id_UNIQUE",
		IndexType: gorn.DBIndexTypeUnique,
		Columns: []*gorn.DBIndexColumn{
			{ColumnName: "user_id", ASC: true},
		},
	})
}
