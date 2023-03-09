package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/jeessy2/ddns-go/v4/config"
)

//go:embed writing.html
var writingEmbedFile embed.FS

const VersionEnv = "DDNS_GO_VERSION"

type writtingData struct {
	Jsonconf          template.JS
	NotAllowWanAccess bool
	config.User
	config.Webhook
	Version string
}
type configData struct {
	DnsName          string
	DnsID            string
	DnsSecret        string
	TTL              string
	Ipv4Enable       string
	Ipv4GetType      string
	Ipv4Url          string
	Ipv4NetInterface string
	Ipv4Cmd          string
	Ipv4Domains      string
	Ipv6Enable       string
	Ipv6GetType      string
	Ipv6Url          string
	Ipv6NetInterface string
	Ipv6Cmd          string
	IPv6Reg          string
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

	confa, err := config.GetConfigCache()
	if err != nil {
		confa.NotAllowWanAccess = true
	}
	tmpl.Execute(writer, &writtingData{
		Jsonconf:          template.JS(getJson(confa.Dnsconfig)),
		NotAllowWanAccess: confa.NotAllowWanAccess,
		User:              confa.User,
		Webhook:           confa.Webhook,
		Version:           os.Getenv(VersionEnv),
	})
}

func getJson(dnsconf []config.Config) string {
	jsonconf := []configData{}
	for _, conf := range dnsconf {
		// 已存在配置文件，隐藏真实的ID、Secret
		idHide, secretHide := getHideIDSecret(&conf)
		jsonconf = append(jsonconf, configData{
			DnsName:          conf.DNS.Name,
			DnsID:            idHide,
			DnsSecret:        secretHide,
			TTL:              conf.TTL,
			Ipv4Enable:       BooltoOn(conf.Ipv4.Enable),
			Ipv4GetType:      conf.Ipv4.GetType,
			Ipv4Url:          conf.Ipv4.URL,
			Ipv4NetInterface: conf.Ipv4.NetInterface,
			Ipv4Cmd:          conf.Ipv4.Cmd,
			Ipv4Domains:      strings.Join(conf.Ipv4.Domains, "\r\n"),
			Ipv6Enable:       BooltoOn(conf.Ipv6.Enable),
			Ipv6GetType:      conf.Ipv6.GetType,
			Ipv6Url:          conf.Ipv6.URL,
			Ipv6NetInterface: conf.Ipv6.NetInterface,
			Ipv6Cmd:          conf.Ipv6.Cmd,
			IPv6Reg:          conf.Ipv6.IPv6Reg,
			Ipv6Domains:      strings.Join(conf.Ipv6.Domains, "\r\n"),
		})
	}
	byt, _ := json.Marshal(jsonconf)
	return string(byt)
}

// 显示的数量
const displayCount int = 3

// hideIDSecret 隐藏真实的ID、Secret
func getHideIDSecret(conf *config.Config) (idHide string, secretHide string) {
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

func BooltoOn(b bool) string {
	if b {
		return "on"
	}
	return ""
}
