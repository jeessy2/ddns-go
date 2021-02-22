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
		if ip.IsLoopback() {
			return true
		}

		_, ipNet192, _ := net.ParseCIDR("192.168.0.0/16")
		if ipNet192.Contains(ip) {
			return true
		}

		_, ipNet172, _ := net.ParseCIDR("172.16.0.0/12")
		if ipNet172.Contains(ip) {
			return true
		}

		_, ipNet10, _ := net.ParseCIDR("10.0.0.0/8")
		if ipNet10.Contains(ip) {
			return true
		}

		_, ipNet100, _ := net.ParseCIDR("100.0.0.0/8")
		if ipNet100.Contains(ip) {
			return true
		}

		_, ipNetFE, _ := net.ParseCIDR("fe80::/10")
		if ipNetFE.Contains(ip) {
			return true
		}

		_, ipNetV6FD, _ := net.ParseCIDR("fd00::/8")
		if ipNetV6FD.Contains(ip) {
			return true
		}

	}

	return false
}
