package web

import (
	"encoding/json"
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
	confa, err := config.GetConfigCache()
	firstTime := err != nil
	// 验证安全性后才允许设置保存配置文件：
	// 内网访问或在服务启动的 1 分钟内
	if (!util.IsPrivateNetwork(request.RemoteAddr) || !util.IsPrivateNetwork(request.Host)) &&
		firstTime && time.Now().Unix()-startTime > 60 { // 1 minutes
		writer.Write([]byte("出于安全考虑，若通过公网访问，仅允许在ddns-go启动的 1 分钟内完成首次配置"))
		return
	}

	confa.NotAllowWanAccess = request.FormValue("NotAllowWanAccess") == "on"
	confa.Username = strings.TrimSpace(request.FormValue("Username"))
	confa.Password = request.FormValue("Password")
	confa.WebhookURL = strings.TrimSpace(request.FormValue("WebhookURL"))
	confa.WebhookRequestBody = strings.TrimSpace(request.FormValue("WebhookRequestBody"))
	// 如启用公网访问，帐号密码不能为空
	if !confa.NotAllowWanAccess && (confa.Username == "" || confa.Password == "") {
		writer.Write([]byte("启用外网访问, 必须输入登录用户名/密码"))
		return
	}

	orignal := request.FormValue("Orignal")
	if orignal != getJson(confa.Dnsconfig) {
		writer.Write([]byte("写入冲突"))
		return
	}

	jsonconf := []string{}
	err = json.Unmarshal([]byte(request.FormValue("Jsonconf")), &jsonconf)
	if err != nil {
		writer.Write([]byte("解析失败"))
		return
	}
	dnsconf := []config.Config{}
	for i, s := range jsonconf {
		if s == "" {
			continue
		}
		v := configData{}
		json.Unmarshal([]byte(s), &v)
		conf := config.Config{TTL: v.TTL}
		// 覆盖以前的配置
		conf.DNS.Name = v.DnsName
		conf.DNS.ID = strings.TrimSpace(v.DnsID)
		conf.DNS.Secret = strings.TrimSpace(v.DnsSecret)
		conf.Ipv4.Enable = v.Ipv4Enable == "on"
		conf.Ipv4.GetType = v.Ipv4GetType
		conf.Ipv4.URL = strings.TrimSpace(v.Ipv4Url)
		conf.Ipv4.NetInterface = v.Ipv4NetInterface
		conf.Ipv4.Cmd = strings.TrimSpace(v.Ipv4Cmd)
		conf.Ipv4.Domains = strings.Split(v.Ipv4Domains, "\r\n")
		conf.Ipv6.Enable = v.Ipv6Enable == "on"
		conf.Ipv6.GetType = v.Ipv6GetType
		conf.Ipv6.URL = strings.TrimSpace(v.Ipv6Url)
		conf.Ipv6.NetInterface = v.Ipv6NetInterface
		conf.Ipv6.Cmd = strings.TrimSpace(v.Ipv6Cmd)
		conf.Ipv6.IPv6Reg = strings.TrimSpace(v.IPv6Reg)
		conf.Ipv6.Domains = strings.Split(v.Ipv6Domains, "\r\n")

		ipCmd := [...]string{"", ""}
		if i < len(confa.Dnsconfig) {
			c := &confa.Dnsconfig[i]
			idHide, secretHide := getHideIDSecret(c)
			if conf.DNS.ID == idHide {
				conf.DNS.ID = c.DNS.ID
			}
			if conf.DNS.Secret == secretHide {
				conf.DNS.Secret = c.DNS.Secret
			}
			ipCmd[0] = c.Ipv4.Cmd
			ipCmd[1] = c.Ipv6.Cmd
		}
		// 修改cmd需要验证：启动前已经保存了帐号密码
		if os.Getenv(SavedPwdOnStartEnv) != "true" &&
			(ipCmd[0] != conf.Ipv4.Cmd || ipCmd[1] != conf.Ipv6.Cmd) {
			writer.Write([]byte("出于安全考虑，修改\"通过命令获取\"要求启动前已配置帐号密码，请配置帐号密码后并重启ddns-go"))
			return
		}
		dnsconf = append(dnsconf, conf)
	}
	confa.Dnsconfig = dnsconf

	// 保存到用户目录
	err = confa.SaveConfig()

	// 只运行一次
	util.ForceCompare = true
	go dns.RunOnce()

	// 回写错误信息
	if err == nil {
		writer.Write([]byte("ok"))
	} else {
		writer.Write([]byte(err.Error()))
	}

}
