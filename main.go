package main

import (
	"ddns-go/config"
	"ddns-go/dns"
	"ddns-go/util"
	"ddns-go/web"
	"embed"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// 监听地址
var listen = flag.String("l", ":9876", "web server listen address")

// 更新频率(秒)
var every = flag.Int("f", 300, "dns update frequency in second")

//go:embed static
var staticEmbededFiles embed.FS

//go:embed favicon.ico
var faviconEmbededFile embed.FS

func main() {
	flag.Parse()

	if _, err := net.ResolveTCPAddr("tcp", *listen); err != nil {
		log.Fatalf("解析监听地址异常，%s", err)
	}

	if util.IsRunInDocker() {
		run()
	} else {
		runAsService()
	}
}

func run() {
	// 启动静态文件服务
	http.Handle("/static/", http.FileServer(http.FS(staticEmbededFiles)))
	http.Handle("/favicon.ico", http.FileServer(http.FS(faviconEmbededFile)))

	http.HandleFunc("/", config.BasicAuth(web.Writing))
	http.HandleFunc("/save", config.BasicAuth(web.Save))
	http.HandleFunc("/logs", config.BasicAuth(web.Logs))
	http.HandleFunc("/ipv4NetInterface", config.BasicAuth(web.Ipv4NetInterfaces))
	http.HandleFunc("/ipv6NetInterface", config.BasicAuth(web.Ipv6NetInterfaces))

	log.Println("监听", *listen, "...")

	// 定时运行
	go dns.RunTimer(time.Duration(*every) * time.Second)
	err := http.ListenAndServe(*listen, nil)

	if err != nil {
		log.Println("启动端口发生异常, 1分钟后自动关闭DOS窗口", err)
		time.Sleep(time.Minute)
		os.Exit(1)
	}
}
