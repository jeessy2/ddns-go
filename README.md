# ddns-go

<a href="https://github.com/jeessy2/ddns-go/releases/latest"><img alt="GitHub release" src="https://img.shields.io/github/release/jeessy2/ddns-go.svg?logo=github&style=flat-square"></a>

自动获得你的公网 IPV4 或 IPV6 地址，并解析到对应的域名服务。

<!-- TOC -->

- [ddns-go](#ddns-go)
  - [特性](#特性)
  - [界面](#界面)
  - [使用](#使用)
    - [直接执行](#直接执行)
    - [Docker](#docker)
  - [使用IPV6](#使用ipv6)
  - [Webhook](#webhook)
  - [开发&自行编译](#开发自行编译)

<!-- /TOC -->

## 特性

- 支持Mac、Windows、Linux系统，支持ARM、x86架构
- 支持的域名服务商 `Alidns(阿里云)` `Dnspod(腾讯云)` `Cloudflare` `华为云`
- 支持接口/网卡获取IP
- 间隔5分钟同步一次
- 支持多个域名同时解析，公司必备
- 支持多级域名
- 网页中配置，简单又方便，可设置登录用户名和密码
- 网页中方便快速查看最近50条日志，不需要跑docker中查看
- 支持webhook

## 界面

![screenshots](https://raw.githubusercontent.com/jeessy2/ddns-go/master/ddns-web.png)

## 使用

### 直接执行

- 下载并解压[https://github.com/jeessy2/ddns-go/releases](https://github.com/jeessy2/ddns-go/releases)
- 双击运行，程序自动打开[http://127.0.0.1:9876](http://127.0.0.1:9876)，修改你的配置，成功
- [可选] 加入到开机启动中，需自行搜索
- [可选] 支持启动带参数 `-l`监听地址 `-f`间隔时间（秒）。如：`./ddns-go -l 127.0.0.1:9876 -f 300`

### Docker

```bash
docker run -d \
  --name ddns-go \
  --restart=always \
  -p 9876:9876 \
  jeessy/ddns-go
```

- 在浏览器中打开`http://主机IP:9876`，修改你的配置，成功
- [可选] docker中默认不支持ipv6，需自行探索如何开启

## 使用IPV6

- 前提：你的电脑或终端能正常获取IPV6，并能正常访问IPV6
- Windows/Mac：推荐 [直接执行](#直接执行)，Windows/Mac桌面版的docker不支持`--net=host`
- 群晖：
  - 套件中心下载docker并打开
  - 注册表中搜索`ddns-go`并下载
  - 映像 -> 选择`jeessy/ddns-go` -> 启动 -> 高级设置 -> 网络中勾选`使用与 Docker Host 相同的网络`，高级设置中勾选`启动自动重新启动`
  - 在浏览器中打开`http://主机IP:9876`，修改你的配置，成功
- Linux的x86或arm架构，如服务器、xx盒子等等，推荐使用`--net=host`模式，简单点

  ```bash
  docker run -d \
    --name ddns-go \
    --restart=always \
    --net=host \
    jeessy/ddns-go
  ```

- 虚拟机中使用有可能正常获取IPV6，但不能正常访问IPV6, 如: `VMware Workstation` `VirtualBox` ...
- [可选] 使用IPV6后，建议设置登录用户名和密码

## Webhook

- 支持webhook, 域名更新成功或不成功时, 会回调填写的URL
- 支持的变量

  |  变量名   | 说明  |
  |  ----  | ----  |
  | #{ipv4Addr}  | 新的IPV4地址 |
  | #{ipv4Result}  | IPV4地址更新结果: `未改变` `失败` `成功`|
  | #{ipv4Domains}  | IPV4的域名，多个以`,`分割 |
  | #{ipv6Addr}  | 新的IPV6地址 |
  | #{ipv6Result}  | IPV6地址更新结果: `未改变` `失败` `成功`|
  | #{ipv6Domains}  | IPV6的域名，多个以`,`分割 |

- RequestBody为空GET请求，不为空POST请求
- Server酱: `https://sc.ftqq.com/[SCKEY].send?text=主人IPv4变了#{ipv4Addr},域名更新结果:#{ipv4Result}`
- 钉钉:
  - 钉钉电脑端 -> 群设置 -> 智能群助手 -> 添加机器人 -> 自定义
  - 只勾选 `自定义关键词`, 输入的关键字必须包含在RequestBody的content中, 如：`你的公网IP变了`
  - URL中输入钉钉给你的 `Webhook地址`
  - RequestBody中输入 `{"msgtype": "text","text": {"content": "你的公网IP变了：#{ipv4Addr}，域名更新结果：#{ipv4Result}"}}`

## 开发&自行编译

- 如果喜欢从源代码编译自己的版本，可以使用本项目提供的 Makefile 构建
- 开发:
  - 首先使用 `make init` 安装 bindata
  - 使用 `make dev` 动态加载修改后的 `writing.html`
- 编译:
  - 如修改了html, 务必使用 `make bindata` 生成编译需要的静态文件
  - 使用 `make build` 生成本地编译后的 `ddns-go` 可执行文件
  - 使用 `make build_docker_image` 自行编译 Docker 镜像
