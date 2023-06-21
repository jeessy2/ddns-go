package util

import (
	"context"
	"os"
	"testing"
)

// TestCustomDNSResolver 测试能否通过 DNSServerEnv 值的 DNS 服务器解析域名的 IP。
func TestCustomDNSResolver(t *testing.T) {
	os.Setenv(DNSServerEnv, "1.1.1.1:53")
	_, err := CustomDNSResolver().LookupIP(context.Background(), "ip", "cloudflare.com")
	if err != nil {
		t.Errorf("Failed to lookup IP, err: %v", err)
	}
}
