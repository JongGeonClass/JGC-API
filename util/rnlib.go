package util

import "strings"

// =로 이루어진 길이가 length인 라인을 만듭니다.
func BarLine(length int) string {
	return strings.Repeat("=", length)
}
