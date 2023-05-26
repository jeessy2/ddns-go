package util

import "strings"

// WriteString 使用 strings.Builder 生成字符串并返回 string
// https://pkg.go.dev/strings#Builder
func WriteString(strs ...string) string {
	var b strings.Builder
	for _, str := range strs {
		b.WriteString(str)
	}

	return b.String()
}
