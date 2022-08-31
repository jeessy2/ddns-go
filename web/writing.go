package web

import (
	"embed"
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
	config.Config
	Version string
}

// Writing 填写信息
func Writing(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFS(writingEmbedFile, "writing.html")
	if err != nil {
		fmt.Println("Error happened..")
		fmt.Println(err)
		return
	}

	conf, err := config.GetConfigCache()
	if err == nil {
		// 已存在配置文件，隐藏真实的ID、Secret
		idHide, secretHide := getHideIDSecret(&conf)
		conf.DNS.ID = idHide
		conf.DNS.Secret = secretHide
		tmpl.Execute(writer, &writtingData{Config: conf, Version: os.Getenv(VersionEnv)})
		return
	}

	// 默认值
	if conf.Ipv4.URL == "" {
		conf.Ipv4.URL = "https://myip4.ipip.net, https://ddns.oray.com/checkip, https://ip.3322.net"
		conf.Ipv4.Enable = true
		conf.Ipv4.GetType = "url"
	}
	if conf.Ipv6.URL == "" {
		conf.Ipv6.URL = "https://myip6.ipip.net, https://speed.neu6.edu.cn/getIP.php, https://v6.ident.me"
		conf.Ipv6.GetType = "url"
	}
	if conf.DNS.Name == "" {
		conf.DNS.Name = "alidns"
	}
	// 默认禁止外部访问
	conf.NotAllowWanAccess = true

	tmpl.Execute(writer, &writtingData{Config: conf, Version: os.Getenv(VersionEnv)})
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
