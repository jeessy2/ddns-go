package util

import (
	"context"
	"net"
)

const DNSServerEnv = "DDNS_GO_DNS_SERVER"

// customDNSResolver 使用 Go 内置 DNS 解析器来解析 server 值的 DNS 服务器。
func customDNSResolver(server string) *net.Resolver {
	if server != "" {
		return &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return net.Dial("udp", server)
			},
		}
	}

	// 如果是 Termux 且未设置 DNS 服务器则设置 DNS 服务器为 1.1.1.1
	if isTermux() && server == "" {
		return customDNSResolver("1.1.1.1:53")
	}

	return &net.Resolver{}
}
