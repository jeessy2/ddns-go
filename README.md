# ddns-go
- 自动获得你的IPV4或IPV6并解析到域名中
- 通过web方式配置
```
go get -u github.com/go-bindata/go-bindata/...
go-bindata -debug -pkg util -o util/staticPagesData.go static/pages/...
go-bindata -pkg static -o static/js_css_data.go -fs -prefix "static/" static/
```