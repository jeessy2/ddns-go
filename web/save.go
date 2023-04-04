package web

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jeessy2/ddns-go/v5/config"
	"github.com/jeessy2/ddns-go/v5/dns"
	"github.com/jeessy2/ddns-go/v5/util"
)

var startTime = time.Now().Unix()

const SavedPwdOnStartEnv = "DDNS_GO_SAVED_PWD_ENV"

// Save 保存
func Save(writer http.ResponseWriter, request *http.Request) {
	result := checkAndSave(request)
	dnsConfJsonStr := "[]"
	if result == "ok" {
		conf, _ := config.GetConfigCached()
		dnsConfJsonStr = getDnsConfStr(conf.DnsConf)
	}
	byt, _ := json.Marshal(map[string]string{"result": result, "dnsConf": dnsConfJsonStr})

	writer.Write(byt)
}

func checkAndSave(request *http.Request) string {
	conf, err := config.GetConfigCached()
	firstTime := err != nil
	// 验证安全性后才允许设置保存配置文件：
	// 内网访问或在服务启动的 1 分钟内
	if (!util.IsPrivateNetwork(request.RemoteAddr) || !util.IsPrivateNetwork(request.Host)) &&
		firstTime && time.Now().Unix()-startTime > 60 { // 1 minutes
		return "出于安全考虑，若通过公网访问，仅允许在ddns-go启动的 1 分钟内完成首次配置"
	}

	conf.NotAllowWanAccess = request.FormValue("NotAllowWanAccess") == "on"
	conf.Username = strings.TrimSpace(request.FormValue("Username"))
	conf.Password = request.FormValue("Password")
	conf.WebhookDisable = request.FormValue("WebhookDisable") == "on"
	conf.WebhookURL = strings.TrimSpace(request.FormValue("WebhookURL"))
	conf.WebhookRequestBody = strings.TrimSpace(request.FormValue("WebhookRequestBody"))
	// 如启用公网访问，帐号密码不能为空
	if !conf.NotAllowWanAccess && (conf.Username == "" || conf.Password == "") {
		return "启用外网访问, 必须输入登录用户名/密码"
	}

	dnsConfFromJS := []dnsConf4JS{}
	err = json.Unmarshal([]byte(request.FormValue("DnsConf")), &dnsConfFromJS)
	if err != nil {
		return "解析配置失败，请重试"
	}
	dnsConfArray := []config.DnsConfig{}
	empty := dnsConf4JS{}
	for k, v := range dnsConfFromJS {
		if v == empty {
			continue
		}
		dnsConf := config.DnsConfig{TTL: v.TTL}
		// 覆盖以前的配置
		dnsConf.DNS.Name = v.DnsName
		dnsConf.DNS.ID = strings.TrimSpace(v.DnsID)
		dnsConf.DNS.Secret = strings.TrimSpace(v.DnsSecret)
		dnsConf.Ipv4.Enable = v.Ipv4Enable == "on"
		dnsConf.Ipv4.GetType = v.Ipv4GetType
		dnsConf.Ipv4.URL = strings.TrimSpace(v.Ipv4Url)
		dnsConf.Ipv4.NetInterface = v.Ipv4NetInterface
		dnsConf.Ipv4.Cmd = strings.TrimSpace(v.Ipv4Cmd)
		if strings.Contains(v.Ipv4Domains, "\r\n") {
			dnsConf.Ipv4.Domains = strings.Split(v.Ipv4Domains, "\r\n")
		} else {
			dnsConf.Ipv4.Domains = strings.Split(v.Ipv4Domains, "\n")
		}
		dnsConf.Ipv6.Enable = v.Ipv6Enable == "on"
		dnsConf.Ipv6.GetType = v.Ipv6GetType
		dnsConf.Ipv6.URL = strings.TrimSpace(v.Ipv6Url)
		dnsConf.Ipv6.NetInterface = v.Ipv6NetInterface
		dnsConf.Ipv6.Cmd = strings.TrimSpace(v.Ipv6Cmd)
		dnsConf.Ipv6.IPv6Reg = strings.TrimSpace(v.IPv6Reg)
		if strings.Contains(v.Ipv6Domains, "\r\n") {
			dnsConf.Ipv6.Domains = strings.Split(v.Ipv6Domains, "\r\n")
		} else {
			dnsConf.Ipv6.Domains = strings.Split(v.Ipv6Domains, "\n")
		}

		ipCmd := [...]string{"", ""}
		if k < len(conf.DnsConf) {
			c := &conf.DnsConf[k]
			idHide, secretHide := getHideIDSecret(c)
			if dnsConf.DNS.ID == idHide {
				dnsConf.DNS.ID = c.DNS.ID
			}
			if dnsConf.DNS.Secret == secretHide {
				dnsConf.DNS.Secret = c.DNS.Secret
			}
			ipCmd[0] = c.Ipv4.Cmd
			ipCmd[1] = c.Ipv6.Cmd
		}
		// 修改cmd需要验证：启动前已经保存了帐号密码
		if os.Getenv(SavedPwdOnStartEnv) != "true" &&
			(ipCmd[0] != dnsConf.Ipv4.Cmd || ipCmd[1] != dnsConf.Ipv6.Cmd) {
			return "出于安全考虑，修改\"通过命令获取\"要求启动前已配置帐号密码，请配置帐号密码后并重启ddns-go"
		}
		dnsConfArray = append(dnsConfArray, dnsConf)
	}
	conf.DnsConf = dnsConfArray

	// 保存到用户目录
	err = conf.SaveConfig()

	// 只运行一次
	util.ForceCompareGlobal = true
	go dns.RunOnce()

	// 回写错误信息
	if err != nil {
		return err.Error()
	}
	return "ok"
}
