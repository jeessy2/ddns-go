package util

import (
	"net"
	"strings"
)

// IsPrivateNetwork 是否为私有地址
// https://en.wikipedia.org/wiki/Private_network
func IsPrivateNetwork(remoteAddr string) bool {
	lastIndex := strings.LastIndex(remoteAddr, ":")
	if lastIndex < 1 {
		return false
	}

	remoteAddr = remoteAddr[:lastIndex]

	// ipv6
	if strings.HasPrefix(remoteAddr, "[") && strings.HasSuffix(remoteAddr, "]") {
		remoteAddr = remoteAddr[1 : len(remoteAddr)-1]
	}

	if ip := net.ParseIP(remoteAddr); ip != nil {
		return ip.IsLoopback() || // 127/8, ::1
			ip.IsPrivate() || // 10/8, 172.16/12, 192.168/16, fc00::/7
			ip.IsLinkLocalUnicast() // 169.254/16, fe80::/10
	}

	// localhost
	if remoteAddr == "localhost" {
		return true
	}

	return false
}
