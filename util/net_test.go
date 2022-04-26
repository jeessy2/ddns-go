package util

import (
	"testing"
)

// TestIsPrivateNetwork 测试是否为私有地址
func TestIsPrivateNetwork(t *testing.T) {

	data := map[string]bool{
		"127.0.0.1:9876":    true,
		"[::1]:9876":        true,
		"192.168.1.18:9876": true,
		"172.16.1.18:9876":  true,
		"10.1.1.18:9876":    true,
		"[fe80::1]:9876":    true,
		"[fd00::1]:9876":    true,
		"100.0.0.1:9876":    false,
		"[2409::1]:9876":    false,
		"223.5.5.5:9876":    false,
	}

	for key, value := range data {
		if IsPrivateNetwork(key) != value {
			t.Errorf("%s 校验失败\n", key)
		}

	}
}
