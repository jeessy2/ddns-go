package util

import (
	"context"
	"testing"
)

// TestCustomDNSResolver 测试能否通过 DNSServerEnv 值的 DNS 服务器解析域名的 IP。
func TestCustomDNSResolver(t *testing.T) {
	_, err := customDNSResolver("1.1.1.1:53").LookupIP(context.Background(), "ip", "cloudflare.com")
	if err != nil {
		t.Errorf("Failed to lookup IP, err: %v", err)
	}
}
