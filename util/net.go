package util

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// IsPrivateNetwork 是否为私有地址
// https://en.wikipedia.org/wiki/Private_network
func IsPrivateNetwork(remoteAddr string) bool {
	// removing optional port from remoteAddr
	if strings.HasPrefix(remoteAddr, "[") { // ipv6
		if index := strings.LastIndex(remoteAddr, "]"); index != -1 {
			remoteAddr = remoteAddr[1:index]
		} else {
			return false
		}
	} else { // ipv4
		if index := strings.LastIndex(remoteAddr, ":"); index != -1 {
			remoteAddr = remoteAddr[:index]
		}
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

// get real IP from request
func GetRealIP(r *http.Request) (addr, extra string) {
	addr = r.RemoteAddr
	extra = fmt.Sprintf("Remote:%s; ", addr)
	real := r.Header.Get("X-Real-IP")
	forward := r.Header.Get("X-Forwarded-For")
	forward, _, _ = strings.Cut(forward, ",")
	if real != "" && real != addr {
		addr = real
		extra += fmt.Sprintf("RealIP:%s; ", real)
	}
	if forward != "" && forward != addr {
		addr = forward
		extra += fmt.Sprintf("Forwarded:%s; ", forward)
	}
	return
}
