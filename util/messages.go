package util

import (
	"log"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var logLang = language.English
var logPrinter = message.NewPrinter(logLang)

func init() {

	message.SetString(language.English, "可使用 .\\ddns-go.exe -s install 安装服务运行", "You can use '.\\ddns-go.exe -s install' to install service")
	message.SetString(language.English, "可使用 sudo ./ddns-go -s install 安装服务运行", "You can use 'sudo ./ddns-go -s install' to install service")
	message.SetString(language.English, "监听 %s", "Listen on %s")
	message.SetString(language.English, "配置文件已保存在: %s", "Config file has been saved to: %s")

	message.SetString(language.English, "你的IP %s 没有变化, 域名 %s", "Your's IP %s has not changed! Domain: %s")
	message.SetString(language.English, "新增域名解析 %s 成功! IP: %s", "Added domain %s successfully! IP: %s")
	message.SetString(language.English, "新增域名解析 %s 失败! 异常信息: %s", "Added domain %s failed! Result: %s")

	message.SetString(language.English, "更新域名解析 %s 成功! IP: %s", "Updated domain %s successfully! IP: %s")
	message.SetString(language.English, "更新域名解析 %s 失败! 异常信息: %s", "Updated domain %s failed! Result: %s")

	message.SetString(language.English, "你的IPv4未变化, 未触发 %s 请求", "Your's IPv4 has not changed, %s request has not been triggered")
	message.SetString(language.English, "你的IPv6未变化, 未触发 %s 请求", "Your's IPv6 has not changed, %s request has not been triggered")
	message.SetString(language.English, "Namecheap 不支持更新 IPv6", "Namecheap don't supports IPv6")

	message.SetString(language.English, "dynadot仅支持单域名配置，多个域名请添加更多配置", "dynadot only supports single domain configuration, please add more configurations")

	// http_util
	message.SetString(language.English, "异常信息: %s", "Exception: %s")
	message.SetString(language.English, "查询域名信息发生异常! %s", "Query domain info failed! %s")
	message.SetString(language.English, "返回内容: %s ,返回状态码: %d", "Response body: %s ,Response status code: %d")
	message.SetString(language.English, "通过接口获取IPv4失败! 接口地址: %s", "Get IPv4 from %s failed")
	message.SetString(language.English, "通过接口获取IPv6失败! 接口地址: %s", "Get IPv6 from %s failed")
	message.SetString(language.English, "将不会触发Webhook, 仅在第 3 次失败时触发一次Webhook, 当前失败次数：%d", "Webhook will not be triggered, only trigger once when the third failure, current failure times: %d")
	message.SetString(language.English, "在DNS服务商中未找到根域名: %s", "Root domain not found in DNS provider: %s")

	// webhook
	message.SetString(language.English, "Webhook配置中的URL不正确", "Webhook url is incorrect")
	message.SetString(language.English, "Webhook中的 RequestBody JSON 无效", "Webhook RequestBody JSON is invalid")
	message.SetString(language.English, "Webhook调用成功! 返回数据：%s", "Webhook called successfully! Response body: %s")
	message.SetString(language.English, "Webhook调用失败! 异常信息：%s", "Webhook called failed! Exception: %s")
	message.SetString(language.English, "Webhook Header不正确: %s", "Webhook header is invalid: %s")
	message.SetString(language.English, "请输入Webhook的URL", "Please enter the Webhook url")

	// callback
	message.SetString(language.English, "Callback的URL不正确", "Callback url is incorrect")
	message.SetString(language.English, "Callback调用成功, 域名: %s, IP: %s, 返回数据: %s", "Webhook called successfully! Domain: %s, IP: %s, Response body: %s")
	message.SetString(language.English, "Callback调用失败, 异常信息: %s", "Webhook called failed! Exception: %s")

	// save
	message.SetString(language.English, "必须输入用户名/密码", "Username/Password is required")
	message.SetString(language.English, "密码不安全！尝试使用更复杂的密码", "Password is not secure! Try using a more complex password")
	message.SetString(language.English, "数据解析失败, 请刷新页面重试", "Data parsing failed, please refresh the page and try again")
	message.SetString(language.English, "第 %s 个配置未填写域名", "The %s config does not fill in the domain")

	// config
	message.SetString(language.English, "从网卡获得IPv4失败", "Get IPv4 from network card failed")
	message.SetString(language.English, "从网卡中获得IPv4失败! 网卡名: %s", "Get IPv4 from network card failed! Network card name: %s")
	message.SetString(language.English, "获取IPv4结果失败! 接口: %s ,返回值: %s", "Get IPv4 result failed! Interface: %s ,Result: %s")
	message.SetString(language.English, "获取%s结果失败! 未能成功执行命令：%s, 错误：%q, 退出状态码：%s", "Get %s result failed! Command: %s, Error: %q, Exit status code: %s")
	message.SetString(language.English, "获取%s结果失败! 命令: %s, 标准输出: %q", "Get %s result failed! Command: %s, Stdout: %q")
	message.SetString(language.English, "从网卡获得IPv6失败", "Get IPv6 from network card failed")
	message.SetString(language.English, "从网卡中获得IPv6失败! 网卡名: %s", "Get IPv6 from network card failed! Network card name: %s")
	message.SetString(language.English, "获取IPv6结果失败! 接口: %s ,返回值: %s", "Get IPv6 result failed! Interface: %s ,Result: %s")
	message.SetString(language.English, "未找到第 %d 个IPv6地址! 将使用第一个IPv6地址", "%dth IPv6 address not found! Will use the first IPv6 address")
	message.SetString(language.English, "IPv6匹配表达式 %s 不正确! 最小从1开始", "IPv6 match expression %s is incorrect! Minimum start from 1")
	message.SetString(language.English, "IPv6将使用正则表达式 %s 进行匹配", "IPv6 will use regular expression %s for matching")
	message.SetString(language.English, "匹配成功! 匹配到地址: %s", "Match successfully! Matched address: %s")
	message.SetString(language.English, "没有匹配到任何一个IPv6地址, 将使用第一个地址", "No IPv6 address matched, will use the first address")
	message.SetString(language.English, "未能获取IPv4地址, 将不会更新", "Failed to get IPv4 address, will not update")
	message.SetString(language.English, "未能获取IPv6地址, 将不会更新", "Failed to get IPv6 address, will not update")

	// domains
	message.SetString(language.English, "域名: %s 不正确", "The domain %s is incorrect")
	message.SetString(language.English, "域名: %s 解析失败", "The domain %s resolution failed")
	message.SetString(language.English, "IPv6未改变, 将等待 %d 次后与DNS服务商进行比对", "IPv6 has not changed, will wait %d times to compare with DNS provider")
	message.SetString(language.English, "IPv4未改变, 将等待 %d 次后与DNS服务商进行比对", "IPv4 has not changed, will wait %d times to compare with DNS provider")

	message.SetString(language.English, "本机DNS异常! 将默认使用 %s, 可参考文档通过 -dns 自定义 DNS 服务器", "Local DNS exception! Will use %s by default, you can use -dns to customize DNS server")
	message.SetString(language.English, "等待网络连接: %s", "Waiting for network connection: %s")
	message.SetString(language.English, "%s 后重试...", "Retry after %s")
	message.SetString(language.English, "网络已连接", "The network is connected")

	// main
	message.SetString(language.English, "监听端口发生异常, 请检查端口是否被占用! %s", "Listen port failed, please check if the port is occupied! %s")
	message.SetString(language.English, "Docker中运行, 请在浏览器中打开 http://docker主机IP:9876 进行配置", "Running in Docker, please open http://docker-host-ip:9876 in the browser for configuration")
	message.SetString(language.English, "ddns-go 服务卸载成功", "ddns-go service uninstalled successfully")
	message.SetString(language.English, "ddns-go 服务卸载失败, 异常信息: %s", "ddns-go service uninstalled failed, Exception: %s")
	message.SetString(language.English, "安装 ddns-go 服务成功! 请打开浏览器并进行配置", "Install ddns-go service successfully! Please open the browser and configure it")
	message.SetString(language.English, "安装 ddns-go 服务失败, 异常信息: %s", "Install ddns-go service failed, Exception: %s")
	message.SetString(language.English, "ddns-go 服务已安装, 无需再次安装", "ddns-go service has been installed, no need to install again")
	message.SetString(language.English, "重启 ddns-go 服务成功", "restart ddns-go service successfully")
	message.SetString(language.English, "启动 ddns-go 服务成功", "start ddns-go service successfully")
	message.SetString(language.English, "ddns-go 服务未安装, 请先安装服务", "ddns-go service is not installed, please install the service first")

	// webhook通知
	message.SetString(language.English, "未改变", "no changed")
	message.SetString(language.English, "失败", "failed")
	message.SetString(language.English, "成功", "success")

	// Login
	message.SetString(language.English, "%q 配置文件为空, 超过3小时禁止从公网访问", "%q configuration file is empty, public network access is prohibited for more than 3 hours")
	message.SetString(language.English, "%q 被禁止从公网访问", "%q is prohibited from accessing the public network")
	message.SetString(language.English, "%q 帐号密码不正确", "%q username or password is incorrect")
	message.SetString(language.English, "%q 登录成功", "%q login successfully")
	message.SetString(language.English, "用户名或密码错误", "Username or password is incorrect")
	message.SetString(language.English, "登录失败次数过多，请等待 %d 分钟后再试", "Too many login failures, please try again after %d minutes")
	message.SetString(language.English, "用户名 %s 的密码已重置成功! 请重启ddns-go", "The password of username %s has been reset successfully! Please restart ddns-go")
	message.SetString(language.English, "需在 %s 之前完成用户名密码设置,请重启ddns-go", "Need to complete the username and password setting before %s, please restart ddns-go")

}

func Log(key string, args ...interface{}) {
	log.Println(LogStr(key, args...))
}

func LogStr(key string, args ...interface{}) string {
	return logPrinter.Sprintf(key, args...)
}

func InitLogLang(lang string) string {
	newLang := language.English
	if strings.HasPrefix(lang, "zh") {
		newLang = language.Chinese
	}
	if newLang != logLang {
		logLang = newLang
		logPrinter = message.NewPrinter(logLang)
	}
	return logLang.String()
}
