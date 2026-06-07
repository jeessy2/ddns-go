package util

import (
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

	return false
}

func ClientIPFromRequest(r *http.Request, trustedProxies []string) (string, bool) {
	peerIP := parseIPFromAddr(r.RemoteAddr)
	if peerIP == nil {
		return "", false
	}

	if !isTrustedProxy(peerIP, trustedProxies) {
		return peerIP.String(), true
	}

	if realIP := parseIPFromAddr(r.Header.Get("X-Real-IP")); realIP != nil {
		return realIP.String(), true
	}

	for _, forwardedIP := range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		if ip := parseIPFromAddr(forwardedIP); ip != nil {
			return ip.String(), true
		}
	}

	return "", false
}

func isTrustedProxy(peerIP net.IP, trustedProxies []string) bool {
	for _, trustedProxy := range trustedProxies {
		trustedProxy = strings.TrimSpace(trustedProxy)
		if trustedProxy == "" {
			continue
		}
		if trustedIP := parseIPFromAddr(trustedProxy); trustedIP != nil && trustedIP.Equal(peerIP) {
			return true
		}
		if _, trustedNet, err := net.ParseCIDR(trustedProxy); err == nil && trustedNet.Contains(peerIP) {
			return true
		}
	}
	return false
}

func parseIPFromAddr(addr string) net.IP {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return nil
	}

	if host, _, err := net.SplitHostPort(addr); err == nil {
		addr = host
	}
	addr = strings.TrimPrefix(strings.TrimSuffix(addr, "]"), "[")
	return net.ParseIP(addr)
}

// GetRequestIPStr get IP string from request
func GetRequestIPStr(r *http.Request) (addr string) {
	addr = "Remote: " + r.RemoteAddr
	if r.Header.Get("X-Real-IP") != "" {
		addr = addr + " ,Real-IP: " + r.Header.Get("X-Real-IP")
	}
	if r.Header.Get("X-Forwarded-For") != "" {
		addr = addr + " ,Forwarded-For: " + r.Header.Get("X-Forwarded-For")
	}
	return addr
}
