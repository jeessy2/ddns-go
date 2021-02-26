package main

import (
	"ddns-go/config"
	"ddns-go/util"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
	run()
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

// 以服务方式运行
func runAsService() {
	svcConfig := &service.Config{
		Name:        "ddns-go",
		DisplayName: "ddns-go",
		Description: "简单好用的DDNS。自动更新域名解析到公网IP(支持阿里云、腾讯云dnspod、Cloudflare、华为云)",
		Arguments:   []string{"-l", *listen, "-f", strconv.Itoa(*every)},
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatalln(err)
	}

	// 处理卸载
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "uninstall":
			if err = s.Uninstall(); err == nil {
				log.Println("ddns-go 服务卸载成功!")
			} else {
				log.Printf("ddns-go 服务卸载失败, ERR: %s", err)
			}
			return
		}
	}

	status, err := s.Status()
	if err != nil && status == service.StatusUnknown {
		// 服务未知，创建服务
		if err = s.Install(); err == nil {
			s.Start()
			openExplorer()
			log.Println("安装 ddns-go 服务成功! 程序会一直运行, 包括重启后。")
			log.Println("如需卸载 ddns-go, 使用 sudo ./ddns-go uninstall")
			log.Println("请在浏览器中进行配置。1分钟后自动关闭DOS窗口!")
			time.Sleep(time.Minute)
			return
		}

		log.Printf("安装 ddns-go 服务失败, ERR: %s", err)
	}

	// 正常启动
	s.Run()

}

func openExplorer() {
	_, err := config.GetConfigCache()
	// 未找到配置文件&&不是在docker中运行 才打开浏览器
	if err != nil && !util.IsRunInDocker() {
		addr, err := net.ResolveTCPAddr("tcp", *listen)
		if err != nil {
			return
		}
		url := ""
		if addr.IP.IsGlobalUnicast() {
			url = fmt.Sprintf("http://%s", addr.String())
		} else {
			url = fmt.Sprintf("http://127.0.0.1:%d", addr.Port)
		}
		go util.OpenExplorer(url)
	}
}
