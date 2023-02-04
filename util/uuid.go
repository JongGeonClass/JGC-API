package util

import "github.com/google/uuid"

// 랜덤으로 uuid를 생성해 반환합니다.
func NewUuid() string {
	return uuid.New().String()
}
