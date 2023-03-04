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
	ConfigMap    map[string]writeConfig
	ConfigGlobal config.ConfigGlobal
	Version      string
}
type writeConfig struct {
	writeLabel
	config.Config
}
type writeLabel struct {
	LabelName   string
	LabelID     string
	LabelSecret string
	LabelHelp   template.HTML
}

var serverList = map[string]writeLabel{
	"alidns": {
		LabelName:   "Alidns(阿里云)",
		LabelID:     "AccessKey ID",
		LabelSecret: "AccessKey Secret",
		LabelHelp:   "<a target='_blank' href='https://ram.console.aliyun.com/manage/ak?spm=5176.12818093.nav-right.dak.488716d0mHaMgg'>创建 AccessKey</a>",
	},
	"dnspod": {
		LabelName:   "Dnspod(腾讯云)",
		LabelID:     "ID",
		LabelSecret: "Token",
		LabelHelp:   "<a target='_blank' href='https://console.dnspod.cn/account/token/token'>创建密钥</a>",
	},
	"cloudflare": {
		LabelName:   "Cloudflare",
		LabelID:     "",
		LabelSecret: "Token",
		LabelHelp:   "<a target='_blank' href='https://dash.cloudflare.com/profile/api-tokens'>创建令牌->编辑区域 DNS (使用模板)</a>",
	},
	"huaweicloud": {
		LabelName:   "华为云",
		LabelID:     "Access Key Id",
		LabelSecret: "Secret Access Key",
		LabelHelp:   "<a target='_blank' href='https://console.huaweicloud.com/iam/?locale=zh-cn#/mine/accessKey'>新增访问密钥</a>",
	},
	"callback": {
		LabelName:   "Callback",
		LabelID:     "URL",
		LabelSecret: "RequestBody",
		LabelHelp:   "<a target='_blank' href='https://github.com/jeessy2/ddns-go#callback'>自定义回调</a> 支持的变量 #{ip}, #{domain}, #{recordType}, #{ttl}",
	},
	"baiducloud": {
		LabelName:   "百度云",
		LabelID:     "AccessKey ID",
		LabelSecret: "AccessKey Secre",
		LabelHelp:   "<a target='_blank' href='https://console.bce.baidu.com/iam/?_=1651763238057#/iam/accesslist'>创建 AccessKey</a><br /><a target='_blank' href='https://ticket.bce.baidu.com/#/ticket/create~productId=60&questionId=393&channel=2'>申请工单</a> DDNS 需调用 API ，而百度云相关 API 仅对申请用户开放，使用前请先提交工单申请。",
	},
	"porkbun": {
		LabelName:   "porkbun",
		LabelID:     "API Key",
		LabelSecret: "Secret Key",
		LabelHelp:   "<a target='_blank' href='https://porkbun.com/account/api'>创建 Access</a>",
	},
	"godaddy": {
		LabelName:   "GoDaddy",
		LabelID:     "Key",
		LabelSecret: "Secret",
		LabelHelp:   "<a target='_blank' href='https://developer.godaddy.com/keys'>创建 API KEY</a>",
	},
	"googledomain": {
		LabelName:   "Google Domain",
		LabelID:     "Username",
		LabelSecret: "Password",
		LabelHelp:   "<a target='_blank' href='https://support.google.com/domains/answer/6147083?hl=zh-Hans'>新建动态域名解析记录</a>",
	},
}

// Writing 填写信息
func Writing(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFS(writingEmbedFile, "writing.html")
	if err != nil {
		fmt.Println("Error happened..")
		fmt.Println(err)
		return
	}

	cglobal, err := config.GetConfigGlobal()
	// 默认禁止外部访问
	if err != nil {
		cglobal.NotAllowWanAccess = true
	}
	cmap := config.GetConfigMap()
	wmap := make(map[string]writeConfig, len(serverList))

	for name, label := range serverList {
		conf, set := cmap[name]
		if set && name != "callback" {
			// 已存在配置文件，隐藏真实的ID、Secret
			idHide, secretHide := getHideIDSecret(&conf)
			conf.DNS.ID = idHide
			conf.DNS.Secret = secretHide
		} else if !set {
			// 默认值
			conf.Ipv4.URL = "https://myip4.ipip.net, https://ddns.oray.com/checkip, https://ip.3322.net, https://4.ipw.cn"
			conf.Ipv4.Enable = true
			conf.Ipv4.GetType = "url"
			conf.Ipv6.URL = "https://speed.neu6.edu.cn/getIP.php, https://v6.ident.me, https://6.ipw.cn"
			conf.Ipv6.Enable = true
			conf.Ipv6.GetType = "netInterface"
		}
		wmap[name] = writeConfig{label, conf}
	}

	tmpl.Execute(writer, &writtingData{ConfigMap: wmap, ConfigGlobal: cglobal, Version: os.Getenv(VersionEnv)})
}

// 显示的数量
const displayCount int = 3

// hideIDSecret 隐藏真实的ID、Secret
func getHideIDSecret(conf *config.Config) (idHide string, secretHide string) {
	if len(conf.DNS.ID) > displayCount {
		idHide = conf.DNS.ID[:displayCount] + strings.Repeat("*", len(conf.DNS.ID)-displayCount)
	} else {
		idHide = conf.DNS.ID
	}
	if len(conf.DNS.Secret) > displayCount {
		secretHide = conf.DNS.Secret[:displayCount] + strings.Repeat("*", len(conf.DNS.Secret)-displayCount)
	} else {
		secretHide = conf.DNS.Secret
	}
	return
}
