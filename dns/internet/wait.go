// Package internet implements utilities for checking the Internet connection.
package internet

import (
	"strings"
	"time"

	"github.com/jeessy2/ddns-go/v6/util"
)

const (
	// fallbackDNS used when a fallback occurs.
	fallbackDNS = "1.1.1.1"

	// delay is the delay time for each DNS lookup.
	delay = time.Second * 5
)

// Wait blocks until the Internet is connected.
//
// See also:
//
//   - https://stackoverflow.com/a/50058255
//   - https://github.com/ddev/ddev/blob/v1.22.7/pkg/globalconfig/global_config.go#L776
func Wait(addresses []string) {
	// fallbase in case loopback DNS is unavailable and only once.
	fallback := false

	for {
		for _, addr := range addresses {
			err := util.LookupHost(addr)
			// Internet is connected.
			if err == nil {
				return
			}

			if !fallback && isLoopback(err) {
				util.Log("本机DNS异常! 将默认使用 %s, 可参考文档通过 -dns 自定义 DNS 服务器", fallbackDNS)
				util.SetDNS(fallbackDNS)

				fallback = true
				continue
			}

			util.Log("等待网络连接: %s", err)

			util.Log("%s 后重试...", delay)
			time.Sleep(delay)
		}
	}
}

// isLoopback checks if the error is a loopback error.
func isLoopback(e error) bool {
	return strings.Contains(e.Error(), "[::1]:53")
}
