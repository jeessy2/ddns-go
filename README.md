<a href="https://github.com/jeessy2/ddns-go/releases/latest"><img alt="GitHub release" src="https://img.shields.io/github/release/jeessy2/ddns-go.svg?logo=github&style=flat-square"></a>

# ddns-go
- 自动获得你的公网IPV4或IPV6并解析到域名中
- 支持Mac、Windows、Linux系统，支持ARM、x86架构
- 支持的域名服务商 `Alidns(阿里云)` `Dnspod(腾讯云)` `Cloudflare` `华为云` `Webhook`
- 间隔5分钟同步一次
- 支持多个域名同时解析，公司必备
- 支持多级域名
- 网页中配置，简单又方便，可设置登录用户名和密码
- 网页中方便快速查看最近50条日志，不需要跑docker中查看

## 系统中使用
- 下载并解压[https://github.com/jeessy2/ddns-go/releases](https://github.com/jeessy2/ddns-go/releases)
- 双击运行，程序自动打开[http://127.0.0.1:9876](http://127.0.0.1:9876)，修改你的配置，成功
- [可选] 加入到开机启动中，需自行搜索

## Docker中使用
```
docker run -d \
  --name ddns-go \
  --restart=always \
  -p 9876:9876 \
  jeessy/ddns-go
```
- 在网页中打开`http://主机IP:9876`，修改你的配置，成功
- [可选] docker中默认不支持ipv6，需自行探索如何开启

## 使用IPV6
- 前提：你的电脑或终端能正常获取IPV6，并能正常访问IPV6
- Windows/Mac：推荐在 `系统中使用`，Windows/Mac桌面版的docker不支持`--net=host`
- 群晖：使用docker。1、注册表中搜索`ddns-go`并下载。 2、映像 -> 选择`jeessy/ddns-go` -> 启动 -> 高级设置 -> 网络中勾选`使用与 Docker Host 相同的网络`)
- Linux的x86或arm架构，如服务器、xx盒子等等，推荐使用`--net=host`模式，简单点
  ```
  docker run -d \
    --name ddns-go \
    --restart=always \
    --net=host \
    jeessy/ddns-go
  ```
- 虚拟机有可能正常获取IPV6，但不能正常访问IPV6
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
- 例(URL):  `https://sc.ftqq.com/[SCKEY].send?text=主人IPv4变了#{ipv4Addr},更新结果:#{ipv4Result}`
- 例(RequestBody): `{"text":"你的IPv4已变为#{ipv4Addr}","desp":"更新结果: #{ipv4Result}"}}`


![avatar](https://raw.githubusercontent.com/jeessy2/ddns-go/master/ddns-web.png)

## Development
```
go get -u github.com/go-bindata/go-bindata/...
go-bindata -debug -pkg util -o util/staticPagesData.go static/pages/...
go-bindata -pkg static -o static/js_css_data.go -fs -prefix "static/" static/
```

## Release
```
go-bindata -pkg util -o util/staticPagesData.go static/pages/...
go-bindata -pkg static -o static/js_css_data.go -fs -prefix "static/" static/

# 自动发布
git tag v0.0.x -m "xxx" 
git push --tags
```