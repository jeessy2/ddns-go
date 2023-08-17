package util

import (
	"context"
	"testing"
)

// TestNewDialerResolver 测试传递 DNS 服务器地址时能否设置 dialer.Resolver。
func TestNewDialerResolver(t *testing.T) {
	dialer.Resolver = nil
	NewDialerResolver("1.1.1.1:53")

	if dialer.Resolver == nil {
		t.Error("Failed to set dialer.Resolver")
	}
}

// TestNewNetResolver 测试能否通过 newNetResolver 返回的 net.Resolver 解析域名的 IP。
func TestNewNetResolver(t *testing.T) {
	_, err := newNetResolver("1.1.1.1:53").LookupIP(context.Background(), "ip", "cloudflare.com")
	if err != nil {
		t.Errorf("Failed to lookup IP, err: %v", err)
	}
}
