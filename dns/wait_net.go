package dns

import (
	"log"
	"net/http"
	"time"
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
			resp, err := http.Get(addr)
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
