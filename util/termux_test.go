package util

import (
	"os"
	"testing"
)

// TestIsTermux 测试在或不在 Termux 中运行都能正确判断
func TestIsTermux(t *testing.T) {
	// 模拟在 Termux 中运行
	os.Setenv("PREFIX", "/data/data/com.termux/files/usr")

	if !isTermux() {
		t.Error("期待 isTermux 返回 true，但得到 false。")
	}

	// 清除 PREFIX 变量，模拟不在 Termux 中运行
	os.Unsetenv("PREFIX")

	if isTermux() {
		t.Error("期待 isTermux 返回 false，但得到 true。")
	}
}
