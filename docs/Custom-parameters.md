# 自定义参数

支持传递自定义参数，目前只实现如下DNS服务商:

## Aliyun
- 通过 [Line参数](https://help.aliyun.com/document_detail/29807.html) 修改：
```
example.com?Line=telecom
```
- 通过 RecordId 修改：
```
example.com?RecordId=xx
```
- 可搭配更多[参数](https://help.aliyun.com/document_detail/29774.html)
```
example.com?RecordId=xx&Line=oversea&Lang=en
```

## Dnspod
- 通过 [record_line参数](https://docs.dnspod.cn/dns/dns-record-line/) 修改：
```
example.com?record_line=境内
```
- 通过 record_id 修改：
```
example.com?record_id=xx
```
- 可搭配更多[参数](https://docs.dnspod.cn/api/modify-records/)
```
example.com?record_id=xx&record_line=境内&status=disable
```
## Callback
- `URL` 或 `RequestBody` 参考, myid 取 domains 中的 myid：
```
https://mycallback.com?ip=#{ip}&domain=#{domain}&myid=#{myid}
```
- domains
```
example.com?myid=xx
```
## Cloudflare
proxied (v4.2.2)
```
example.com?proxied=true
```

更多请见 [#336](https://github.com/jeessy2/ddns-go/issues/336)