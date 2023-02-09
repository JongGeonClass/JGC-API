package dbmodel

import "github.com/thak1411/gorn"

// JGC에서 사용할 인덱스들입니다.
var indexes []*gorn.DBIndex = nil

// 새로운 인덱스를 등록합니다.
func AddIndex(index *gorn.DBIndex) {
	indexes = append(indexes, index)
}

// 등록된 모든 인덱스를 가져옵니다.
func GetIndexes() []*gorn.DBIndex {
	return indexes
}
