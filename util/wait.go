package util

import (
	"strings"
	"time"
)

const (
	// delay is the delay time for each DNS lookup.
	delay = time.Second * 5
)

// Wait blocks until the Internet is connected.
//
// See also:
//
//   - https://stackoverflow.com/a/50058255
//   - https://github.com/ddev/ddev/blob/v1.22.7/pkg/globalconfig/global_config.go#L776
func WaitInternet(addresses []string, fallbackDNS []string) {
	// fallbase in case loopback DNS is unavailable and only once.
	fallbackTimes := 0

	for {
		for _, addr := range addresses {

			err := LookupHost(addr)
			// Internet is connected.
			if err == nil {
				return
			}

			Log("等待网络连接: %s", err)
			Log("%s 后重试...", delay)

			if isLoopbackErr(err) && fallbackTimes >= 10 {
				dns := fallbackDNS[fallbackTimes%len(fallbackDNS)]
				Log("DNS异常! 将默认使用 %s, 可参考文档通过 -dns 自定义 DNS 服务器", dns)
				SetDNS(dns)
			}

			fallbackTimes = fallbackTimes + 1
			time.Sleep(delay)
		}
	}
}

// isLoopbackErr checks if the error is a loopback error.
func isLoopbackErr(e error) bool {
	return strings.Contains(e.Error(), "[::1]:53")
}
