// based on https://github.com/v2fly/v2ray-core/blob/0c5abc7e53aed41480dd2a07f611bd34c753e880/transport/internet/system_dns_android_test.go

package util

import (
	"context"
	"os"
	"testing"
)

// TestCustomDNSResolver 测试能否通过 DNSServerEnv 值的 DNS 服务器解析域名的 IP。
func TestCustomDNSResolver(t *testing.T) {
	os.Setenv(DNSServerEnv, "223.5.5.5:53")
	_, err := customDNSResolver().LookupIP(context.Background(), "ip", "aliyun.com")
	if err != nil {
		t.Errorf("Failed to lookup IP, err: %v", err)
	}
}
