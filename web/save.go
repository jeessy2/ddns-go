package web

import (
	"ddns-go/config"
	"ddns-go/dns"
	"ddns-go/util"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"
)

// Save 保存
func Save(writer http.ResponseWriter, request *http.Request) {

	conf := &config.Config{}
	// 初始化以前的配置
	conf.InitConfigFromFile()

	idNew := request.FormValue("DnsID")
	secretNew := request.FormValue("DnsSecret")

	idHide, secretHide := getHideIDSecret(conf)

	if idNew != idHide {
		conf.DNS.ID = idNew
	}
	if secretNew != secretHide {
		conf.DNS.Secret = secretNew
	}

	// 覆盖以前的配置
	conf.DNS.Name = request.FormValue("DnsName")

	conf.Ipv4.Enable = request.FormValue("Ipv4Enable") == "on"
	conf.Ipv4.URL = request.FormValue("Ipv4Url")
	conf.Ipv4.Domains = strings.Split(request.FormValue("Ipv4Domains"), "\r\n")

	conf.Ipv6.Enable = request.FormValue("Ipv6Enable") == "on"
	conf.Ipv6.URL = request.FormValue("Ipv6Url")
	conf.Ipv6.Domains = strings.Split(request.FormValue("Ipv6Domains"), "\r\n")

	// 保存到用户目录
	util.GetConfigFilePath()
	byt, err := yaml.Marshal(conf)
	if err != nil {
		log.Println(err)
	}

	ioutil.WriteFile(util.GetConfigFilePath(), byt, 0600)

	// 只运行一次
	go dns.RunOnce()

	// 跳转
	http.Redirect(writer, request, "/?saveOk=true", http.StatusFound)

}
