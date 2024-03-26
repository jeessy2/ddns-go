package util

import (
	"math/rand"
	"strings"
	"time"
)

// Wait blocks until the Internet is connected.
//
// See also:
//
//   - https://stackoverflow.com/a/50058255
//   - https://github.com/ddev/ddev/blob/v1.22.7/pkg/globalconfig/global_config.go#L776
func WaitInternet(addresses []string) {
	delay := time.Second * 5

	for {
		for _, addr := range addresses {

			err := LookupHost(addr)
			// Internet is connected.
			if err == nil {
				return
			}

			Log("等待网络连接: %s", err)
			Log("%s 后重试...", delay)

			if isDNSErr(err) && len(DNSList) > 0 {
				dns := DNSList[rand.Intn(len(DNSList))]
				Log("本机DNS异常! 将默认使用 %s, 可参考文档通过 -dns 自定义 DNS 服务器", dns)
				SetDNS(dns)
			}

			time.Sleep(delay)
		}
	}
}

// isDNSErr checks if the error is caused by DNS.
func isDNSErr(e error) bool {
	return strings.Contains(e.Error(), "[::1]:53: read: connection refused")
}
