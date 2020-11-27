package main

import (
	"ddns-go/config"
	"ddns-go/dns"
	"ddns-go/static"
	"ddns-go/util"
	"ddns-go/web"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

const port = "9876"

func main() {
	listen := flag.String("l", ":9876", "web server listen address")
	every := flag.String("f", "300", "dns update frequency in second")
	flag.Parse()
	// 启动静态文件服务
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(static.AssetFile())))
	http.Handle("/favicon.ico", http.StripPrefix("/", http.FileServer(static.AssetFile())))

	http.HandleFunc("/", config.BasicAuth(web.Writing))
	http.HandleFunc("/save", config.BasicAuth(web.Save))
	http.HandleFunc("/logs", config.BasicAuth(web.Logs))

	addr, err := net.ResolveTCPAddr("tcp", *listen)
	if err != nil {
		log.Fatalf("解析监听地址异常，%s", err)
	}
	url := ""
	if addr.IP.IsGlobalUnicast() {
		url = fmt.Sprintf("http://%s", addr.String())
	} else if addr.IP.To4() != nil || addr.IP == nil || addr.IP.Equal(net.ParseIP("::")) {
		url = fmt.Sprintf("http://127.0.0.1:%d", addr.Port)
	} else {
		url = fmt.Sprintf("http://[::1]:%d", addr.Port)
	}
	// 未找到配置文件才打开浏览器
	_, err = config.GetConfigCache()
	if err != nil {
		go util.OpenExplorer(url)
	}

	log.Println("监听", *listen, "...")

	// 定时运行
	delay, err := strconv.ParseUint(*every, 10, 64)
	if err != nil {
		delay = 300
	}
	go dns.RunTimer(time.Duration(delay) * time.Second)
	err = http.ListenAndServe(*listen, nil)

	if err != nil {
		log.Println("启动端口发生异常, 1分钟后自动关闭此窗口", err)
		time.Sleep(time.Minute)
	}

}
