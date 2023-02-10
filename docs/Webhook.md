# Webhook

- 支持 Webhook, 域名更新成功或不成功时, 会回调填写的 `URL`
- 支持的变量：

  |  变量名   | 说明  |
  |  ----  | ----  |
  | #{ipv4Addr}  | 新的IPv4地址 |
  | #{ipv4Result}  | IPv4地址更新结果: `未改变` `失败` `成功`|
  | #{ipv4Domains}  | IPv4的域名，多个以`,`分割 |
  | #{ipv6Addr}  | 新的IPv6地址 |
  | #{ipv6Result}  | IPv6地址更新结果: `未改变` `失败` `成功`|
  | #{ipv6Domains}  | IPv6的域名，多个以`,`分割 |

- `RequestBody` 为空 GET 请求，不为空 POST 请求

## Server酱
```
https://sctapi.ftqq.com/[SendKey].send?title=主人IPv4变了#{ipv4Addr},域名更新结果:#{ipv4Result}
```

## Bark
```
https://api.day.app/[YOUR_KEY]/主人IPv4变了#{ipv4Addr},域名更新结果:#{ipv4Result}
```

## 钉钉
1. 钉钉电脑端 -> 群设置 -> 智能群助手 -> 添加机器人 -> 自定义
2. 只勾选 `自定义关键词`, 输入的关键字必须包含在RequestBody的content中, 如：`你的公网IP变了`
3. `URL` 中输入钉钉给你的 `Webhook URL`
4. `RequestBody` 中输入：
```JSON
{
	"msgtype": "markdown",
	"markdown": {
		"title": "你的公网IP变了",
		"text": "#### 你的公网IP变了 \n - IPV4地址：#{ipv4Addr} \n - 域名更新结果：#{ipv4Result} \n"
	}
}
```

## 飞书
1. 飞书电脑端 -> 群设置 -> 添加机器人 -> 自定义机器人
2. 安全设置只勾选 `自定义关键词`, 输入的关键字必须包含在 `RequestBody` 的 content 中, 如：`你的公网IP变了`
3. `URL` 中输入飞书给你的 `Webhook URL`
4. `RequestBody` 中输入：
```JSON
{
	"msg_type": "post",
	"content": {
		"post": {
			"zh_cn": {
				"title": "你的公网IP变了",
				"content": [
					[{
						"tag": "text",
						"text": "IPV4地址：#{ipv4Addr}"
					}],
					[{
						"tag": "text",
						"text": "域名更新结果：#{ipv4Result}"
					}]
				]
			}
		}
	}
}
```

## Telegram
### [ddns-telegram-bot](https://github.com/WingLim/ddns-telegram-bot)
1. `Webhook URL`：`https://api.telegram.org/bot{your token}/sendmessage`
2. `RequestBody`：
```JSON
{
   "chat_id":"{notify room id}",
   "text":"外網IP改變：\n新IPV4地址：#{ipv4Addr}\n域名更新结果：#{ipv4Result}"
}
```
3. Result：

![image](https://user-images.githubusercontent.com/15572804/199481480-a47685a1-cdf1-4f57-9e7d-fae9433a7f8f.png)

## plusplus 推送加
1. [获取 token](https://www.pushplus.plus/push1.html)
2. `URL` 中输入 `https://www.pushplus.plus/send`
3. `RequestBody` 中输入：
```JSON
{
	"token": "your token",
	"title": "你的公网IP变了",
	"content": "你的公网IP变了 \n - IPV4地址：#{ipv4Addr} \n - 域名更新结果：#{ipv4Result} \n"
}
```

## Discord
1. Discord任意客户端 -> 伺服器 -> 频道设置 -> 整合 -> 查看Webhook -> 新Webhook -> 复制Webhook网址
2. `URL` 中输入Discord复制的 `Webhook URL`
3. `RequestBody` 中输入：
```JSON
{
	"content": "域名 #{ipv4Domains} 动态解析 #{ipv4Result}.",
	"embeds": [{
		"description": "#{ipv4Domains} 的动态解析 #{ipv4Result}, IP: #{ipv4Addr}",
		"color": 15258703,
		"author": {
			"name": "DDNS"
		},
		"footer": {
			"text": "DDNS #{ipv4Result}"
		}
	}]
}
```

## 企业微信
1. 下载 企业微信 → 左上角三横杠 → 全新创建企业 → 个人组件团队(创建个人的企业群聊)
2. 进入群聊添加 [群机器人] 复制机器人 Webhook 地址填入 ddns-go 后台 `Webhook URL` 地址栏。
3. 在 `RequestBody` 栏填入回调函数，格式：
```JSON
{
	"msgtype": "markdown",
	"markdown": {
		"content": "公网IP变更：\n 新IPV6地址：#{ipv6Addr} \n 已解析的域名：#{ipv6Domains} \n 域名更新结果：#{ipv6Result}"
	}
}
```

## 饭碗警告
链接：[https://fwalert.com](https://fwalert.com/321991) (含 aff，注册后会赠送 10 元余额)

支持通过 邮件(0.02元/次)、短信(0.1元/次)、电话(0.2元/次)、Telegram(免费)、饭碗警告App(免费) 的方式推送通知。

1. 先点击右上角头像选择“联系方式”并在此对你所希望的通知渠道进行绑定，然后进入“转发规则”，点击加号新建规则，其中触发方式选 `Webhook`，可参照下图添加模板变量，并依据你所设置的模板变量来设置通知正文，最后联系方式选择先前绑定的通知渠道即可。
![Snipaste_2022-07-29_10-32-35](https://user-images.githubusercontent.com/51308700/181670740-cb0c2a9a-6250-430a-a5d9-77d7fa796e45.png)
![Snipaste_2022-07-29_10-36-27](https://user-images.githubusercontent.com/51308700/181671132-8595a9b0-34b3-4bcc-9d52-3e48285246ee.png)
2. 保存转发规则后会生成一个 `Webhook` 地址，将该地址后添加 `?result=#{ipv6Result}&addr=#{ipv6Addr}` (此处等号前的变量需与前面设置的一致) 填入 ddns-go 后台 `Webhook URL` 地址栏并保持 `RequestBody` 留空即可。

## Apprise
Apprise 的邮箱推送
1. apprise Webhook URL  
`https://你的公网域名:端口/notify/你的密钥` 就是一个 `Webhook URL`   
"你的密钥"是自建的 `{key}` 或者 `token`，可以换成任意一个
2. 设置推送邮箱  
进入 `https://你的公网域名:端口/cfg/你的密钥` 配置一下  
`mailto://邮箱账号:授权码@qq.com?name=📢DDNS-GO`
3. 配置 ddns-go  
在 `Webhook URL` 中填入第一步里的 `URL`  
在 `RequestBody` 中填入：
```JSON
{
    "title": "公网IP变动了",
    "format": "html",
    "body": "新IPV4地址：#{ipv4Addr}\n已解析的域名：#{ipv4Domains}\n域名更新结果：#{ipv4Result}\n\n------DDNS-GO------"
}
```
*其中的 "title"、"type"、"format"、"body" 都是 apprise 定义的，其中的 #{ipv4Addr}、#{ipv4Domains}、#{ipv4Result} 是 ddns-go 定义的

效果如图：
![微信图片_20220920090907](https://user-images.githubusercontent.com/17892238/191145478-12343a0f-8183-4a62-80b1-b513e8e83ed5.jpg)

## ntfy
[ntfy](https://ntfy.sh/) : 免费免注册可自建的多平台推送方案。
- 使用官方/公共服务器，推荐以 uuid 为 topic ：  
1. `uuidgen` : `e056a473-c080-4f34-b49c-e62f9fcd1f9d`  
2. `URL` ：`https://ntfy.sh/`  
3. `RequestBody` ：
```JSON
{
    "topic": "e056a473-c080-4f34-b49c-e62f9fcd1f9d",
    "message": "IPv4已变更：#{ipv4Addr}，域名 #{ipv4Domains} 更新#{ipv4Result}",
    "title": "DDNS-GO Cloudflare 更新",
    "priority": 2,
    "actions": [{ "action": "view", "label": "管理界面", "url": "http://192.168.0.1:9876/" }]
}
```
4. 客户端添加订阅 topic：`e056a473-c080-4f34-b49c-e62f9fcd1f9d` ，可设置别名。

- 自建服务并且设置了认证：  
1. 生成 `auth` 参数(*nix命令)：  
`echo -n "Basic `echo -n '\<user>:\<pass>' | base64`" | base64 | tr -d '='`  
（替换 `<user>` 和 `<pass>`），请结合 `https`加密 `URL`，详细请参考 [ntfy文档](https://docs.ntfy.sh/publish/#query-param)。  
2. URL： `https://ntfy.example.com/?auth=<上一步生成的base64 auth参数>`  
3. RequestBody ：
```JSON
{
    "topic": "ddns-go",
    "message": "IPv4已变更：#{ipv4Addr}，域名 #{ipv4Domains} 更新#{ipv4Result}",
    "title": "DDNS-GO Cloudflare 更新",
    "priority": 2,
    "actions": [{ "action": "view", "label": "管理界面", "url": "http://192.168.0.1:9876/" }]
}
```
4. 客户端在设置里更改默认服务器为自建：`https://ntfy.example.com/`，设置用户名和密码，然后添加订阅 topic：`ddns-go` 。
- 推送效果

![ddnsgo-ntfy](https://user-images.githubusercontent.com/86276507/208280040-c9483679-4b22-4c82-83fd-865990f120fd.png)

## Gotify
1. 首先，登录到 Gotify 的 WebUI，点击 `APPS` -> `CREATE APPLICATION` 来创建 Token，得到 Token 后回到 `ddns-go`。
2. 然后，登录到 ddns-go，找到 `Webhook`，在 `URL` 处填入：  
`http://[IP]/message?token=[Token]`  
将 [IP] 替换为 Gotify 服务器的 IP，将 [Token] 替换为得到的 Token。
在 `RequestBody` 处填入：
```JSON
{
	"title": "你的公网 IP 变了",
	"message": "IPv4 地址：#{ipv4Addr}\n域名更新结果：#{ipv4Result}",
        "priority": 5,
	"extras": {
		"client::display": {
			"contentType": "text/plain"
		}
	}
}
```
效果：

![result](https://user-images.githubusercontent.com/62788816/216801381-c2b89896-eb78-4c30-aa43-60304c76b8d8.png)

参考：
1. [Push messages · Gotify](https://gotify.net/docs/pushmsg)
2. [Message Extras · Gotify](https://gotify.net/docs/msgextras)

更多请见 [#327](https://github.com/jeessy2/ddns-go/issues/327)