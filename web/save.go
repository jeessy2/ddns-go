package web

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/jeessy2/ddns-go/v5/config"
	"github.com/jeessy2/ddns-go/v5/dns"
	"github.com/jeessy2/ddns-go/v5/util"
)

var startTime = time.Now().Unix()

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
	usernameNew := strings.TrimSpace(request.FormValue("Username"))
	passwordNew := request.FormValue("Password")

	// 验证安全性后才允许设置保存配置文件：
	if time.Now().Unix()-startTime > 5*60 {
		firstTime := err != nil

		// 首次设置 && 通过外网访问 必需在服务启动的 5 分钟内
		if firstTime &&
			(!util.IsPrivateNetwork(request.RemoteAddr) || !util.IsPrivateNetwork(request.Host)) {
			return "若通过公网访问，仅允许在ddns-go启动后 5 分钟内完成首次配置"
		}

		// 非首次设置 && 从未设置过帐号密码 && 本次设置了帐号或密码 必须在5分钟内
		if !firstTime &&
			(conf.Username == "" && conf.Password == "") &&
			(usernameNew != "" || passwordNew != "") {
			return "若从未设置过帐号密码，仅允许在ddns-go启动后 5 分钟内设置，请重启ddns-go"
		}

	}

	// 如果密码不为空则检查是否够强
	if passwordNew != "" {
		err = validate(passwordNew)
		if err != nil {
			return err.Error()
		}
	}

	conf.NotAllowWanAccess = request.FormValue("NotAllowWanAccess") == "on"
	conf.Username = usernameNew
	conf.Password = passwordNew
	conf.WebhookURL = strings.TrimSpace(request.FormValue("WebhookURL"))
	conf.WebhookRequestBody = strings.TrimSpace(request.FormValue("WebhookRequestBody"))
	conf.WebhookHeaders = strings.TrimSpace(request.FormValue("WebhookHeaders"))

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

		if k < len(conf.DnsConf) {
			c := &conf.DnsConf[k]
			idHide, secretHide := getHideIDSecret(c)
			if dnsConf.DNS.ID == idHide {
				dnsConf.DNS.ID = c.DNS.ID
			}
			if dnsConf.DNS.Secret == secretHide {
				dnsConf.DNS.Secret = c.DNS.Secret
			}

			// 修改cmd需要验证：必须设置帐号密码
			if (conf.Username == "" && conf.Password == "") &&
				(c.Ipv4.Cmd != dnsConf.Ipv4.Cmd || c.Ipv6.Cmd != dnsConf.Ipv6.Cmd) {
				return "修改 \"通过命令获取\" 必须设置帐号密码，请先设置帐号密码"
			}
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
