<a href="https://github.com/jeessy2/ddns-go/releases/latest"><img alt="GitHub release" src="https://img.shields.io/github/release/jeessy2/ddns-go.svg?logo=github&style=flat-square"></a>

# ddns-go
- 自动获得你的公网IPV4或IPV6并解析到域名中
- 支持Mac、Windows、Linux系统，支持ARM、x86架构
- 间隔5分钟同步一次
- 支持多个域名同时解析，公司必备
- 支持多级域名
- 支持的域名服务商 `Alidns(阿里云)` `Dnspod(腾讯云)` `Cloudflare` 

## 系统中使用
- 下载[https://github.com/jeessy2/ddns-go/releases](https://github.com/jeessy2/ddns-go/releases)
- 双击运行，程序自动打开[http://127.0.0.1:9876](http://127.0.0.1:9876)，修改你的配置，成功

## Docker中使用
```
docker run -d \
  --name ddns-go \
  --restart=always \
  -p 127.0.0.1:9876:9876 \
  jeessy/ddns-go
```
- 在docker主机上打开[http://127.0.0.1:9876](http://127.0.0.1:9876)，修改你的配置，成功

![avatar](ddns-web.png)

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