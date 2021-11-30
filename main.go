package main

import (
	"ddns-go/config"
	"ddns-go/dns"
	"ddns-go/util"
	"ddns-go/web"
	"embed"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kardianos/service"
)

// 监听地址
var listen = flag.String("l", ":9876", "监听地址")

// 更新频率(秒)
var every = flag.Int("f", 300, "同步间隔时间(秒)")

// 服务管理
var serviceType = flag.String("s", "", "服务管理, 支持install, uninstall")

// 配置文件路径
var configFilePath = flag.String("c", "", "自定义配置文件路径")

//go:embed static
var staticEmbededFiles embed.FS

//go:embed favicon.ico
var faviconEmbededFile embed.FS

func main() {
	flag.Parse()
	if _, err := net.ResolveTCPAddr("tcp", *listen); err != nil {
		log.Fatalf("解析监听地址异常，%s", err)
	}

	if *configFilePath != "" {
		absPath, _ := filepath.Abs(*configFilePath)
		os.Setenv(util.ConfigFilePathENV, absPath)
	}
	switch *serviceType {
	case "install":
		installService()
	case "uninstall":
		uninstallService()
	default:
		if util.IsRunInDocker() {
			run(100 * time.Millisecond)
		} else {
			s := getService()
			status, _ := s.Status()
			if status != service.StatusUnknown {
				// 以服务方式运行
				s.Run()
			} else {
				// 非服务方式运行
				switch s.Platform() {
				case "windows-service":
					log.Println("可使用 .\\ddns-go.exe -s install 安装服务运行")
				default:
					log.Println("可使用 ./ddns-go -s install 安装服务运行")
				}
				run(100 * time.Millisecond)
			}
		}
	}
}

func run(firstDelay time.Duration) {
	// 启动静态文件服务
	http.Handle("/static/", http.FileServer(http.FS(staticEmbededFiles)))
	http.Handle("/favicon.ico", http.FileServer(http.FS(faviconEmbededFile)))

	http.HandleFunc("/", web.BasicAuth(web.Writing))
	http.HandleFunc("/save", web.BasicAuth(web.Save))
	http.HandleFunc("/logs", web.BasicAuth(web.Logs))
	http.HandleFunc("/clearLog", web.BasicAuth(web.ClearLog))
	http.HandleFunc("/ipv4NetInterface", web.BasicAuth(web.Ipv4NetInterfaces))
	http.HandleFunc("/ipv6NetInterface", web.BasicAuth(web.Ipv6NetInterfaces))
	http.HandleFunc("/webhookTest", web.BasicAuth(web.WebhookTest))

	log.Println("监听", *listen, "...")

	// 没有配置, 自动打开浏览器
	autoOpenExplorer()

	// 定时运行
	go dns.RunTimer(firstDelay, time.Duration(*every)*time.Second)
	err := http.ListenAndServe(*listen, nil)

	if err != nil {
		log.Println("启动端口发生异常, 请检查端口是否被占用", err)
		time.Sleep(time.Minute)
		os.Exit(1)
	}
}

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	// 服务运行，延时10秒运行，等待网络
	run(10 * time.Second)
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func getService() service.Service {
	nowUsers, _ := user.Current()
	options := make(service.KeyValue)
	_, isOpenWRT := os.Stat("/sbin/procd")
	if nowUsers.Name == "root" {
		log.Println("正在以管理员模式执行服务...")
		options["UserService"] = false
	} else {
		log.Println("正在以用户模式执行服务...")
		options["UserService"] = true
	}
	if isOpenWRT == nil {
	options["SysvScript"] = `#!/bin/sh /etc/rc.common
NAME={{.DisplayName}}
DESCRIPTION="{{.Description}}"
CMD="{{.Path}}"
USE_PROCD=1
START=99

start_service() {
	echo Starting $NAME service...
	echo 正在启动 $NAME 服务...
	procd_open_instance $NAME
	procd_set_param respawn
	procd_set_param command $CMD
	procd_append_param command {{range .Arguments}} {{.|cmd}}{{end}}
	procd_set_param stdout 1
	procd_close_instance
}

stop_service() {
	echo Stopping $NAME service...
	echo 正在停止 $NAME 服务...
}

reload_service() {
	procd_send_signal clash
}

restart() {
	stop
	echo
	start
}
`
	} else {
		options["SysvScript"] = `#!/bin/sh
# For RedHat and cousins:
# chkconfig: - 99 01
# description: {{.Description}}
# processname: {{.Path}}

### BEGIN INIT INFO
# Provides:          {{.Path}}
# Required-Start:
# Required-Stop:
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: {{.DisplayName}}
# Description:       {{.Description}}
### END INIT INFO

cmd="{{.Path}}{{range .Arguments}} {{.|cmd}}{{end}}"

name=$(basename $(readlink -f $0))
pid_file="/var/run/$name.pid"
stdout_log="/var/log/$name.log"
stderr_log="/var/log/$name.err"

[ -e /etc/sysconfig/$name ] && . /etc/sysconfig/$name

get_pid() {
	cat "$pid_file"
}

is_running() {
	[ -f "$pid_file" ] && ps $(get_pid) >/dev/null 2>&1
}

case "$1" in
start)
	if is_running; then
		echo "$name Already started"
		echo "$name 已启动"
	else
		echo "Starting $name service..."
		echo "正在启动 $name 服务..."
		{{if .WorkingDirectory}}cd '{{.WorkingDirectory}}'{{end}}
		$cmd >>"$stdout_log" 2>>"$stderr_log" &
		echo $! >"$pid_file"
		if ! is_running; then
			echo "Unable to start, see $stdout_log and $stderr_log"
			echo "该服务无法启用， 详细信息请见 $stdout_log 和 $stderr_log"
			exit 1
		fi
	fi
	;;
stop)
	if is_running; then
		echo -n "Stopping $name service..."
		echo -n "正在停止 $name 服务..."
		kill $(get_pid)
		for i in $(seq 1 10); do
			if ! is_running; then
				break
			fi
			echo -n "."
			sleep 1
		done
		echo
		if is_running; then
			echo "Not stopped; may still be shutting down or shutdown may have failed"
			echo "无法停止；服务可能正在关闭中或关闭向导运行失败"
			exit 1
		else
			echo "$name Stopped"
			echo "$name 已停止"
			if [ -f "$pid_file" ]; then
				rm "$pid_file"
			fi
		fi
	else
		echo "$name Not running"
		echo "$name 未运行"
	fi
	;;
restart)
	$0 stop
	if is_running; then
		echo "Unable to stop, will not attempt to start"
		echo "无法停止服务，将不会尝试启动"
		exit 1
	fi
	$0 start
	;;
status)
	if is_running; then
		echo "Running"
		echo "运行中"
	else
		echo "Stopped"
		echo "已停止"
		exit 1
	fi
	;;
*)
	echo "Usage: $0 {start|stop|restart|status}"
	echo "用法： $0 {start|stop|restart|status}"
	exit 1
	;;
esac
exit 0

`
	}
	svcConfig := &service.Config {
			Name:        "ddns-go",
			DisplayName: "ddns-go",
			Description: "简单好用的DDNS。自动更新域名解析到公网IP(支持阿里云、腾讯云dnspod、Cloudflare、华为云)",
			Arguments:   []string{"-l", *listen, "-f", strconv.Itoa(*every), "-c", *configFilePath},
			Option:      options,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatalln(err)
	}
	return s
}

// 卸载服务
func uninstallService() {
	s := getService()

	status, _ := s.Status()
	// 处理卸载
	if status != service.StatusUnknown {
		s.Stop()
		if err := s.Uninstall(); err == nil {
			log.Println("ddns-go 服务卸载成功!")
		} else {
			log.Printf("ddns-go 服务卸载失败, ERR: %s\n", err)
		}
	} else {
		log.Printf("ddns-go 服务未安装")
	}
}

// 安装服务
func installService() {
	s := getService()

	status, err := s.Status()
	if err != nil && status == service.StatusUnknown {
		// 服务未知，创建服务
		if err = s.Install(); err == nil {
			s.Start()
			log.Println("安装 ddns-go 服务成功! 程序会一直运行, 包括重启后。")
			log.Println("请在浏览器进行配置! 如果不存在配置文件, 会自动打开浏览器。")
			return
		}

		log.Printf("安装 ddns-go 服务失败, ERR: %s\n", err)
		switch s.Platform() {
		case "windows-service":
			log.Println("请确保使用如下命令: .\\ddns-go.exe -s install")
		default:
			log.Println("请确保使用如下命令: ./ddns-go -s install")
		}
	}

	if status != service.StatusUnknown {
		log.Println("ddns-go 服务已安装, 无需在次安装")
	}
}

// 打开浏览器
func autoOpenExplorer() {
	_, err := config.GetConfigCache()
	// 未找到配置文件
	if err != nil {
		if util.IsRunInDocker() {
			// docker中运行, 提示
			fmt.Println("Docker中运行, 请在浏览器中打开 http://docker主机IP:端口 进行配置")
		} else {
			// 主机运行, 打开浏览器
			addr, err := net.ResolveTCPAddr("tcp", *listen)
			if err != nil {
				return
			}
			url := fmt.Sprintf("http://127.0.0.1:%d", addr.Port)
			if addr.IP.IsGlobalUnicast() {
				url = fmt.Sprintf("http://%s", addr.String())
			}
			go util.OpenExplorer(url)
		}
	}
}
