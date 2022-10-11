# ddns-go

[![GitHub release](https://img.shields.io/github/release/jeessy2/ddns-go.svg?logo=github&style=flat-square)](https://github.com/jeessy2/ddns-go/releases/latest) [![](https://goreportcard.com/badge/github.com/jeessy2/ddns-go/v4)](https://goreportcard.com/report/github.com/jeessy2/ddns-go/v4) [![](https://img.shields.io/docker/image-size/jeessy/ddns-go)](https://registry.hub.docker.com/r/jeessy/ddns-go) [![](https://img.shields.io/docker/pulls/jeessy/ddns-go)](https://registry.hub.docker.com/r/jeessy/ddns-go)

自动获得你的公网 IPv4 或 IPv6 地址，并解析到对应的域名服务。

<!-- TOC -->

- [ddns-go](#ddns-go)
  - [特性](#特性)
  - [系统中使用](#系统中使用)
  - [Docker中使用](#docker中使用)
  - [使用IPv6](#使用ipv6)
  - [Webhook](#webhook)
  - [Callback](#callback)
  - [界面](#界面)
  - [开发&自行编译](#开发自行编译)

<!-- /TOC -->

## 特性

- 支持Mac、Windows、Linux系统，支持ARM、x86架构
- 支持的域名服务商 `Alidns(阿里云)` `Dnspod(腾讯云)` `Cloudflare` `华为云` `Callback` `百度云` `porkbun` `GoDaddy` `Google Domains`
- 支持接口/网卡获取IP
- 支持以服务的方式运行
- 默认间隔5分钟同步一次
- 支持多个域名同时解析，公司必备
- 支持多级域名
- 网页中配置，简单又方便，可设置 `登录用户名和密码` / `禁止从公网访问`
- 网页中方便快速查看最近50条日志，不需要跑docker中查看
- 支持webhook通知
- 支持TTL
- 支持部分dns服务商传递自定义参数，实现地域解析等功能

## 系统中使用

- 下载并解压[https://github.com/jeessy2/ddns-go/releases](https://github.com/jeessy2/ddns-go/releases)
- 双击运行, 如没有找到配置, 程序自动打开[http://127.0.0.1:9876](http://127.0.0.1:9876)
- [可选] 安装服务
  - Mac/Linux: `sudo ./ddns-go -s install`
  - Win(以管理员打开cmd): `.\ddns-go.exe -s install`
  - 安装服务也支持 `-l`监听地址 `-f`同步间隔时间(秒) `-c`自定义配置文件路径 `-noweb`不启动web服务
- [可选] 服务卸载
  - Mac/Linux: `sudo ./ddns-go -s uninstall`
  - Win(以管理员打开cmd): `.\ddns-go.exe -s uninstall`
- [可选] 支持启动带参数 `-l`监听地址 `-f`同步间隔时间(秒) `-c`自定义配置文件路径 `-noweb`不启动web服务。如：`./ddns-go -l 127.0.0.1:9876 -f 600 -c /Users/name/ddns-go.yaml`

## Docker中使用

- 不挂载主机目录, 删除容器同时会删除配置

  ```bash
  # host模式, 同时支持IPv4/IPv6, Liunx系统推荐
  docker run -d --name ddns-go --restart=always --net=host jeessy/ddns-go
  # 桥接模式, 只支持IPv4, Mac/Windows系统推荐
  docker run -d --name ddns-go --restart=always -p 9876:9876 jeessy/ddns-go
  ```

- 在浏览器中打开`http://主机IP:9876`，修改你的配置，成功
- [可选] 挂载主机目录, 删除容器后配置不会丢失。可替换 `/opt/ddns-go` 为主机目录, 配置文件为隐藏文件

  ```bash
  docker run -d --name ddns-go --restart=always --net=host -v /opt/ddns-go:/root jeessy/ddns-go
  ```

- [可选] 支持启动带参数 `-l`监听地址 `-f`间隔时间(秒)

  ```bash
  docker run -d --name ddns-go --restart=always --net=host jeessy/ddns-go -l :9877 -f 600
  ```

## 使用IPv6

- 前提：你的电脑或终端能正常获取IPv6，并能正常访问IPv6
- Windows/Mac：推荐 [系统中使用](#系统中使用)，Windows/Mac桌面版的docker不支持`--net=host`
- 群晖：
  - 套件中心下载docker并打开
  - 注册表中搜索`ddns-go`并下载
  - 映像 -> 选择`jeessy/ddns-go` -> 启动 -> 高级设置 -> 网络中勾选`使用与 Docker Host 相同的网络`，高级设置中勾选`启动自动重新启动`
  - 在浏览器中打开`http://群晖IP:9876`，修改你的配置，成功
- Linux的x86或arm架构，推荐使用Docker的`--net=host`模式。参考 [Docker中使用](#Docker中使用)
- 虚拟机中使用有可能正常获取IPv6，但不能正常访问IPv6
- [可选] 使用IPv6后，建议勾选`禁止从公网访问`

## Webhook

- 支持webhook, 域名更新成功或不成功时, 会回调填写的URL
- 支持的变量

  |  变量名   | 说明  |
  |  ----  | ----  |
  | #{ipv4Addr}  | 新的IPv4地址 |
  | #{ipv4Result}  | IPv4地址更新结果: `未改变` `失败` `成功`|
  | #{ipv4Domains}  | IPv4的域名，多个以`,`分割 |
  | #{ipv6Addr}  | 新的IPv6地址 |
  | #{ipv6Result}  | IPv6地址更新结果: `未改变` `失败` `成功`|
  | #{ipv6Domains}  | IPv6的域名，多个以`,`分割 |

- RequestBody为空GET请求，不为空POST请求
- Server酱: `https://sctapi.ftqq.com/[SendKey].send?title=主人IPv4变了#{ipv4Addr},域名更新结果:#{ipv4Result}`
- Bark: `https://api.day.app/[YOUR_KEY]/主人IPv4变了#{ipv4Addr},域名更新结果:#{ipv4Result}`
- 钉钉:
  - 钉钉电脑端 -> 群设置 -> 智能群助手 -> 添加机器人 -> 自定义
  - 只勾选 `自定义关键词`, 输入的关键字必须包含在RequestBody的content中, 如：`你的公网IP变了`
  - URL中输入钉钉给你的 `Webhook地址`
  - RequestBody中输入 `{"msgtype": "markdown","markdown": {"title":"你的公网IP变了","text": "#### 你的公网IP变了 \n - IPV4地址：#{ipv4Addr} \n - 域名更新结果：#{ipv4Result} \n"}}`
- 飞书：
  - 飞书电脑端 -> 群设置 -> 添加机器人 -> 自定义机器人
  - 安全设置只勾选 `自定义关键词`, 输入的关键字必须包含在RequestBody的content中, 如：`你的公网IP变了`
  - URL中输入飞书给你的 `Webhook地址`
  - RequestBody中输入 `{"msg_type": "post","content": {"post": {"zh_cn": {"title": "你的公网IP变了","content": [[{"tag": "text","text": "IPV4地址：#{ipv4Addr}"}] , [{"tag": "text","text": "域名更新结果：#{ipv4Result}"}]]}}}}`
- Telegram: [ddns-telegram-bot](https://github.com/WingLim/ddns-telegram-bot)

## Callback

- 通过自定义回调可支持更多的第三方DNS服务商
- 配置的域名有几行, 就会回调几次
- 支持的变量

  |  变量名   | 说明  |
  |  ----  | ----  |
  | #{ip}  | 新的IPv4/IPv6地址 |
  | #{domain}  | 当前域名 |
  | #{recordType}  | 记录类型 `A`或`AAAA` |
  | #{ttl}  | ttl |
- RequestBody为空GET请求，不为空POST请求

## 界面

![screenshots](https://raw.githubusercontent.com/jeessy2/ddns-go/master/ddns-web.png)

## 开发&自行编译

- 如果喜欢从源代码编译自己的版本，可以使用本项目提供的 Makefile 构建
- 开发环境 golang 1.19
- 使用 `make build` 生成本地编译后的 `ddns-go` 可执行文件
- 使用 `make build_docker_image` 自行编译 Docker 镜像
