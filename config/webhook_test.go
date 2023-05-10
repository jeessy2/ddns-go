package config

import (
	"fmt"
	"testing"
)

// TestParseHeaderArr 测试 parseHeaderArr
func TestParseHeaderArr(t *testing.T) {
	headers := "a : 1\r\nb:2\r\n"
	expected := `map[a:1 b:2]`
	parsedHeaders := checkParseHeaders(headers)
	resultStr := fmt.Sprintf("%v", parsedHeaders)
	if resultStr != expected {
		t.Error("解析Header失败", resultStr)
	}
}
