package internal

import (
	"log"
	"strings"
	"time"

	"github.com/jeessy2/ddns-go/v5/util"
)

// waitForNetworkConnected 等待网络连接后继续
//
// addresses：用于测试网络是否连接的域名
func WaitForNetworkConnected(addresses []string) {
	// 延时 5 秒
	timeout := time.Second * 5

	loopbackServer := "[::1]:53"
	find := false

	for {
		for _, addr := range addresses {
			// https://github.com/jeessy2/ddns-go/issues/736
			client := util.CreateHTTPClient()
			resp, err := client.Get(addr)
			if err != nil {

				// 如果 err 包含回环地址（[::1]:53）则表示没有 DNS 服务器，设置 DNS 服务器
				if strings.Contains(err.Error(), loopbackServer) && !find {
					server := "1.1.1.1:53"
					log.Printf("解析回环地址 %s 失败！将默认使用 %s，可参考文档通过 -dns 自定义 DNS 服务器",
						loopbackServer, server)

					util.NewDialerResolver(server)
					find = true
					continue
				}

				log.Printf("等待网络连接：%s。%s 后重试...", err, timeout)
				// 等待 5 秒后重试
				time.Sleep(timeout)
				continue
			}

			// 网络已连接
			resp.Body.Close()
			return
		}
	}
}
