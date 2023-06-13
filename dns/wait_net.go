package dns

import (
	"log"
	"time"

	"github.com/jeessy2/ddns-go/v5/util"
)

// waitForNetworkConnected 等待网络连接后继续
func waitForNetworkConnected() {
	// 延时 5 秒
	timeout := time.Second * 5

	// 测试网络是否连接的域名
	addresses := []string{
		alidnsEndpoint,
		baiduEndpoint,
		zonesAPI,
		recordListAPI,
		googleDomainEndpoint,
		huaweicloudEndpoint,
		nameCheapEndpoint,
		porkbunEndpoint,
		tencentCloudEndPoint,
	}

	for {
		for _, addr := range addresses {
			// https://github.com/jeessy2/ddns-go/issues/736
			client := util.CreateHTTPClient()
			resp, err := client.Get(addr)
			if err != nil {
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
