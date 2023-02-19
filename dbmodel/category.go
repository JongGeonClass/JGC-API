package dbmodel

import (
	"encoding/json"
	"fmt"

	"github.com/thak1411/gorn"
)

// GROUP_CONCAT으로 category 리스트를 가져올 때 사용할 객체입니다.
type CategoryList []Category

// DB에서 Scan할 때 struct 로 변환해줍니다.
func (c *CategoryList) Scan(v interface{}) error {
	var pbyte []byte
	switch w := v.(type) {
	case []byte:
		pbyte = w
	case string:
		pbyte = []byte(w)
	case nil:
		return nil
	default:
		return fmt.Errorf("unsupported type: %v", w)
	}
	return json.Unmarshal(pbyte, c)
}

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
