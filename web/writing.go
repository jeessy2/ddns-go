package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/jeessy2/ddns-go/v6/config"
)

//go:embed writing.html
var writingEmbedFile embed.FS

const VersionEnv = "DDNS_GO_VERSION"

// js中的dns配置
type dnsConf4JS struct {
	Name             string
	DnsName          string
	DnsID            string
	DnsSecret        string
	TTL              string
	Ipv4Enable       bool
	Ipv4GetType      string
	Ipv4Url          string
	Ipv4NetInterface string
	Ipv4Cmd          string
	Ipv4Domains      string
	Ipv6Enable       bool
	Ipv6GetType      string
	Ipv6Url          string
	Ipv6NetInterface string
	Ipv6Cmd          string
	Ipv6Reg          string
	Ipv6Domains      string
}

// Writing 填写信息
func Writing(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFS(writingEmbedFile, "writing.html")
	if err != nil {
		fmt.Println("Error happened..")
		fmt.Println(err)
		return
	}

	conf, err := config.GetConfigCached()
	// 默认禁止公网访问
	if err != nil {
		conf.NotAllowWanAccess = true
	}

	ipv4, ipv6, _ := config.GetNetInterface()

	err = tmpl.Execute(writer, struct {
		DnsConf           template.JS
		NotAllowWanAccess bool
		Username          string
		config.Webhook
		Version string
		Ipv4    []config.NetInterface
		Ipv6    []config.NetInterface
	}{
		DnsConf:           template.JS(getDnsConfStr(conf.DnsConf)),
		NotAllowWanAccess: conf.NotAllowWanAccess,
		Username:          conf.User.Username,
		Webhook:           conf.Webhook,
		Version:           os.Getenv(VersionEnv),
		Ipv4:              ipv4,
		Ipv6:              ipv6,
	})
	if err != nil {
		fmt.Println("Error happened..")
		fmt.Println(err)
	}
}

func getDnsConfStr(dnsConf []config.DnsConfig) string {
	dnsConfArray := []dnsConf4JS{}
	for _, conf := range dnsConf {
		// 已存在配置文件，隐藏真实的ID、Secret
		idHide, secretHide := getHideIDSecret(&conf)
		dnsConfArray = append(dnsConfArray, dnsConf4JS{
			Name:             conf.Name,
			DnsName:          conf.DNS.Name,
			DnsID:            idHide,
			DnsSecret:        secretHide,
			TTL:              conf.TTL,
			Ipv4Enable:       conf.Ipv4.Enable,
			Ipv4GetType:      conf.Ipv4.GetType,
			Ipv4Url:          conf.Ipv4.URL,
			Ipv4NetInterface: conf.Ipv4.NetInterface,
			Ipv4Cmd:          conf.Ipv4.Cmd,
			Ipv4Domains:      strings.Join(conf.Ipv4.Domains, "\r\n"),
			Ipv6Enable:       conf.Ipv6.Enable,
			Ipv6GetType:      conf.Ipv6.GetType,
			Ipv6Url:          conf.Ipv6.URL,
			Ipv6NetInterface: conf.Ipv6.NetInterface,
			Ipv6Cmd:          conf.Ipv6.Cmd,
			Ipv6Reg:          conf.Ipv6.Ipv6Reg,
			Ipv6Domains:      strings.Join(conf.Ipv6.Domains, "\r\n"),
		})
	}
	byt, _ := json.Marshal(dnsConfArray)
	return string(byt)
}

// 显示的数量
const displayCount int = 3

// hideIDSecret 隐藏真实的ID、Secret
func getHideIDSecret(conf *config.DnsConfig) (idHide string, secretHide string) {
	if len(conf.DNS.ID) > displayCount && conf.DNS.Name != "callback" {
		idHide = conf.DNS.ID[:displayCount] + strings.Repeat("*", len(conf.DNS.ID)-displayCount)
	} else {
		idHide = conf.DNS.ID
	}
	if len(conf.DNS.Secret) > displayCount && conf.DNS.Name != "callback" {
		secretHide = conf.DNS.Secret[:displayCount] + strings.Repeat("*", len(conf.DNS.Secret)-displayCount)
	} else {
		secretHide = conf.DNS.Secret
	}
	return
}
