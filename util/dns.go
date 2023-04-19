// based on https://github.com/v2fly/v2ray-core/blob/0c5abc7e53aed41480dd2a07f611bd34c753e880/transport/internet/system_dns_android.go

package util

import (
	"context"
	"net"
	"os"
)

const DNSServerEnv = "DDNS_GO_DNS_SERVER"

// customDNSResolver 当 DNSServerEnv 值不为空时，使用 Go 内置 DNS 解析器来解析其 DNS 服务器。
func customDNSResolver() *net.Resolver {
	s := os.Getenv(DNSServerEnv)
	if s != "" {
		return &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial("udp", s)
			},
		}
	}
	return &net.Resolver{}
}
