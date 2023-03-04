package web

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jeessy2/ddns-go/v4/config"
	"github.com/jeessy2/ddns-go/v4/dns"
	"github.com/jeessy2/ddns-go/v4/util"
)

var startTime = time.Now().Unix()

const SavedPwdOnStartEnv = "DDNS_GO_SAVED_PWD_ENV"

// Save 保存
func Save(writer http.ResponseWriter, request *http.Request) {
	cglobal, err := config.GetConfigGlobal()
	firstTime := err != nil
	// 验证安全性后才允许设置保存配置文件：
	// 内网访问或在服务启动的 1 分钟内
	if (!util.IsPrivateNetwork(request.RemoteAddr) || !util.IsPrivateNetwork(request.Host)) &&
		firstTime &&
		time.Now().Unix()-startTime > 60 { // 1 minutes
		writer.Write([]byte("出于安全考虑，若通过公网访问，仅允许在ddns-go启动的 1 分钟内完成首次配置"))
		return
	}

	cglobal.Username = strings.TrimSpace(request.FormValue("Username"))
	cglobal.Password = request.FormValue("Password")
	cglobal.WebhookURL = strings.TrimSpace(request.FormValue("WebhookURL"))
	cglobal.WebhookRequestBody = strings.TrimSpace(request.FormValue("WebhookRequestBody"))
	cglobal.NotAllowWanAccess = request.FormValue("NotAllowWanAccess") == "on"

	// 如启用公网访问，帐号密码不能为空
	if !cglobal.NotAllowWanAccess && (cglobal.Username == "" || cglobal.Password == "") {
		writer.Write([]byte("启用外网访问, 必须输入登录用户名/密码"))
		return
	}

	cmap := config.GetConfigMap()
	smap := make(map[string]config.Config, len(serverList))
	for name := range serverList {
		conf := cmap[name]

		idNew := strings.TrimSpace(request.FormValue(name + "_DnsID"))
		secretNew := strings.TrimSpace(request.FormValue(name + "_DnsSecret"))

		idHide, secretHide := getHideIDSecret(&conf)
		if idNew != idHide {
			conf.DNS.ID = idNew
		}
		if secretNew != secretHide {
			conf.DNS.Secret = secretNew
		}

		ipv4CmdInput := strings.TrimSpace(request.FormValue(name + "_Ipv4Cmd"))
		ipv6CmdInput := strings.TrimSpace(request.FormValue(name + "_Ipv6Cmd"))
		// 修改cmd需要验证：
		// 启动前已经保存了帐号密码
		if os.Getenv(SavedPwdOnStartEnv) != "true" &&
			(ipv4CmdInput != conf.Ipv4.Cmd || ipv6CmdInput != conf.Ipv6.Cmd) {
			writer.Write([]byte("出于安全考虑，修改\"通过命令获取\"要求启动前已配置帐号密码，请配置帐号密码后并重启ddns-go"))
			return
		}

		// 覆盖以前的配置
		conf.Ipv4.Enable = request.FormValue(name+"_Ipv4Enable") == "on"
		conf.Ipv4.URL = strings.TrimSpace(request.FormValue(name + "_Ipv4Url"))
		conf.Ipv4.GetType = request.FormValue(name + "_Ipv4GetType")
		conf.Ipv4.NetInterface = request.FormValue(name + "_Ipv4NetInterface")
		conf.Ipv4.Domains = strings.Split(request.FormValue(name+"_Ipv4Domains"), "\r\n")
		conf.Ipv4.Cmd = ipv4CmdInput

		conf.Ipv6.Enable = request.FormValue(name+"_Ipv6Enable") == "on"
		conf.Ipv6.GetType = request.FormValue(name + "_Ipv6GetType")
		conf.Ipv6.NetInterface = request.FormValue(name + "_Ipv6NetInterface")
		conf.Ipv6.URL = strings.TrimSpace(request.FormValue(name + "_Ipv6Url"))
		conf.Ipv6.IPv6Reg = strings.TrimSpace(request.FormValue(name + "_IPv6Reg"))
		conf.Ipv6.Domains = strings.Split(request.FormValue(name+"_Ipv6Domains"), "\r\n")
		conf.Ipv6.Cmd = ipv6CmdInput

		conf.TTL = request.FormValue(name + "TTL")
		smap[name] = conf
	}

	// 保存到用户目录
	err = config.SaveConfig(cglobal, smap)

	// 只运行一次
	util.ForceCompare = true
	dns.DNSInit()
	go dns.RunOnce()

	// 回写错误信息
	if err == nil {
		writer.Write([]byte("ok"))
	} else {
		writer.Write([]byte(err.Error()))
	}

}
