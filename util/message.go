package util

import (
	"log"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var logPrinter = message.NewPrinter(language.English)

func init() {

	message.SetString(language.English, "可使用 .\\ddns-go.exe -s install 安装服务运行", "You can use 'sudo .\\ddns-go -s install' to install service")
	message.SetString(language.English, "可使用 sudo ./ddns-go -s install 安装服务运行", "You can use 'sudo ./ddns-go -s install' to install service")
	message.SetString(language.English, "监听 %q", "Listen on %q")
	message.SetString(language.English, "配置文件已保存在: %s", "Config file has been saved to: %s")

	message.SetString(language.English, "你的IP %s 没有变化, 域名 %s", "Your's IP %s has not changed! Domain: %s")
	message.SetString(language.English, "新增域名解析 %s 成功! IP: %s", "Added domain %s successfully! IP: %s")
	message.SetString(language.English, "新增域名解析 %s 失败! 异常信息: %s", "Added domain %s failed! Result: %s")

	message.SetString(language.English, "更新域名解析 %s 成功! IP: %s", "Updated domain %s successfully! IP: %s")
	message.SetString(language.English, "更新域名解析 %s 失败! 异常信息: %s", "Updated domain %s failed! Result: %s")

	message.SetString(language.English, "你的IPv4未变化, 未触发 %s 请求", "Your's IPv4 has not changed, %s request has not been triggered")
	message.SetString(language.English, "你的IPv6未变化, 未触发 %s 请求", "Your's IPv6 has not changed, %s request has not been triggered")
	message.SetString(language.English, "Namecheap 不支持更新 IPv6", "Namecheap don't supports IPv6")

	// http_util
	message.SetString(language.English, "请求接口 %q 失败", "Request api %q failed")
	message.SetString(language.English, "异常信息: %s", "Exception: %s")
	message.SetString(language.English, "返回内容: %s ,返回状态码: %d", "Response body: %s ,Response status code: %d")

	message.SetString(language.English, "通过接口获取IPv4失败! 接口地址: %s", "Get IPv4 from %s failed")
	message.SetString(language.English, "通过接口获取IPv6失败! 接口地址: %s", "Get IPv6 from %s failed")

	message.SetString(language.English, "将不会触发Webhook, 仅在第 3 次失败时触发一次Webhook, 当前失败次数：%d", "Webhook will not be triggered, only trigger once when the third failure, current failure times: %d")

}

func Log(key string, args ...interface{}) {
	log.Println(LogStr(key, args...))
}

func LogStr(key string, args ...interface{}) string {
	return logPrinter.Sprintf(key, args...)
}

func InitLogLang(lang string) string {
	logLang := language.English
	if strings.HasPrefix(lang, "zh") {
		logLang = language.Chinese
	}
	logPrinter = message.NewPrinter(logLang)
	return logLang.String()
}
