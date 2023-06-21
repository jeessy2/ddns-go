package util

import (
	"context"
	"net"
	"os"
)

const DNSServerEnv = "DDNS_GO_DNS_SERVER"

// CustomDNSResolver 当 DNSServerEnv 值不为空时，使用 Go 内置 DNS 解析器来解析其 DNS 服务器。
func CustomDNSResolver() *net.Resolver {
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
