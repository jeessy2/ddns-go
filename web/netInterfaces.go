package web

import (
	"encoding/json"
	"net/http"

	"github.com/jeessy2/ddns-go/v4/config"
)

// Ipv4NetInterfaces 获得Ipv4网卡信息
func Ipv4NetInterfaces(writer http.ResponseWriter, request *http.Request) {
	ipv4, _, err := config.GetNetInterface()
	if len(ipv4) > 0 && err == nil {
		byt, err := json.Marshal(ipv4)
		if err == nil {
			writer.Write(byt)
			return
		}
	}
}

// Ipv6NetInterfaces 获得Ipv6网卡信息
func Ipv6NetInterfaces(writer http.ResponseWriter, request *http.Request) {
	_, ipv6, err := config.GetNetInterface()
	if len(ipv6) > 0 && err == nil {
		byt, err := json.Marshal(ipv6)
		if err == nil {
			writer.Write(byt)
			return
		}
	}
}
