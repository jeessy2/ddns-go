package util

import (
	"net/http"
	"testing"
)

// TestIsPrivateNetwork 测试是否为私有地址
func TestIsPrivateNetwork(t *testing.T) {

	data := map[string]bool{
		"127.0.0.1":         true, // listen on default port
		"127.0.0.1:9876":    true,
		"[::1]":             true,
		"[::1]:9876":        true,
		"192.168.1.18:9876": true,
		"172.16.1.18:9876":  true,
		"10.1.1.18:9876":    true,
		"[fe80::1]:9876":    true,
		"[fd00::1]:9876":    true,
		"100.0.0.1":         false,
		"100.0.0.1:9876":    false,
		"[2409::1]":         false,
		"[2409::1]:9876":    false,
		"223.5.5.5:9876":    false,
	}

	for key, value := range data {
		if IsPrivateNetwork(key) != value {
			t.Errorf("%s 校验失败\n", key)
		}

	}
}

// test get request IP string from request
func TestGetRequestIPStr(t *testing.T) {
	req := http.Request{RemoteAddr: "192.168.1.1", Header: http.Header{}}
	req.Header.Set("X-Real-IP", "10.0.0.1")
	req.Header.Set("X-Forwarded-For", "10.0.0.2")
	if GetRequestIPStr(&req) != "Remote: 192.168.1.1 ,Real-IP: 10.0.0.1 ,Forwarded-For: 10.0.0.2" {
		t.Errorf("GetRequestIPStr failed")
	}
}
