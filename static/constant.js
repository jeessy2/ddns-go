const DNS_PROVIDERS = {
  alidns: {
    name: {
      "en": "Aliyun",
      "zh-cn": "阿里云",
    },
    idLabel: "AccessKey ID",
    secretLabel: "AccessKey Secret",
    helpHtml: {
      "en": "<a target='_blank' href='https://ram.console.aliyun.com/manage/ak?spm=5176.12818093.nav-right.dak.488716d0mHaMgg'>Create AccessKey</a>",
      "zh-cn": "<a target='_blank' href='https://ram.console.aliyun.com/manage/ak?spm=5176.12818093.nav-right.dak.488716d0mHaMgg'>创建 AccessKey</a>",
    }
  },
  tencentcloud: {
    name: {
      "en": "Tencent",
      "zh-cn": "腾讯云",
    },
    idLabel: "SecretId",
    secretLabel: "SecretKey",
    helpHtml: {
      "en": "<a target='_blank' href='https://console.dnspod.cn/account/token/apikey'>Create AccessKey</a>",
      "zh-cn": "<a target='_blank' href='https://console.dnspod.cn/account/token/apikey'>创建腾讯云 API 密钥</a>",
    }
  },
  dnspod: {
    name: {
      "en": "DnsPod",
    },
    idLabel: "ID",
    secretLabel: "Token",
    helpHtml: {
      "en": "<a target='_blank' href='https://console.dnspod.cn/account/token/token'>Create Token</a>",
      "zh-cn": "<a target='_blank' href='https://console.dnspod.cn/account/token/token'>创建 DNSPod Token</a>",
    }
  },
  cloudflare: {
    name: {
      "en": "Cloudflare",
    },
    idLabel: "",
    secretLabel: "Token",
    helpHtml: {
      "en": "<a target='_blank' href='https://dash.cloudflare.com/profile/api-tokens'>Create Token -> Edit Zone DNS (Use template)</a>",
      "zh-cn": "<a target='_blank' href='https://dash.cloudflare.com/profile/api-tokens'>创建令牌 -> 编辑区域 DNS (使用模板)</a>",
    }
  },
  huaweicloud: {
    name: {
      "en": "Huawei",
      "zh-cn": "华为云",
    },
    idLabel: "Access Key Id",
    secretLabel: "Secret Access Key",
    helpHtml: {
      "en": "<a target='_blank' href='https://console.huaweicloud.com/iam/?locale=zh-cn#/mine/accessKey'>Create</a>",
      "zh-cn": "<a target='_blank' href='https://console.huaweicloud.com/iam/?locale=zh-cn#/mine/accessKey'>新增访问密钥</a>",
    }
  },
  callback: {
    name: {
      "en": "Callback",
    },
    idLabel: "URL",
    secretLabel: "RequestBody",
    helpHtml: {
      "en": "<a target='_blank' href='https://github.com/jeessy2/ddns-go/blob/master/README_EN.md#callback'>Callback</a> Support variables #{ip}, #{domain}, #{recordType}, #{ttl}",
      "zh-cn": "<a target='_blank' href='https://github.com/jeessy2/ddns-go#callback'>自定义回调</a> 支持的变量 #{ip}, #{domain}, #{recordType}, #{ttl}",
    }
  },
  baiducloud: {
    name: {
      "en": "Baidu",
      "zh-cn": "百度云",
    },
    idLabel: "AccessKey ID",
    secretLabel: "AccessKey Secret",
    helpHtml: {
      "en": "<a target='_blank' href='https://console.bce.baidu.com/iam/?_=1651763238057#/iam/accesslist'>Create AccessKey</a><br /><a target='_blank' href='https://ticket.bce.baidu.com/#/ticket/create~productId=60&questionId=393&channel=2'>Apply for a ticket</a> DDNS needs to call the API, and the related APIs of Baidu Cloud are only open to users who have applied. Please submit a ticket before using it.",
      "zh-cn": "<a target='_blank' href='https://console.bce.baidu.com/iam/?_=1651763238057#/iam/accesslist'>创建 AccessKey</a><br /><a target='_blank' href='https://ticket.bce.baidu.com/#/ticket/create~productId=60&questionId=393&channel=2'>申请工单</a> DDNS 需调用 API ，而百度云相关 API 仅对申请用户开放，使用前请先提交工单申请。",
    }
  },
  porkbun: {
    name: {
      "en": "Porkbun",
    },
    idLabel: "API Key",
    secretLabel: "Secret Key",
    helpHtml: {
      "en": "<a target='_blank' href='https://porkbun.com/account/api'>Create Access</a>",
      "zh-cn": "<a target='_blank' href='https://porkbun.com/account/api'>创建 Access</a>",
    }
  },
  godaddy: {
    name: {
      "en": "GoDaddy",
    },
    idLabel: "Key",
    secretLabel: "Secret",
    helpHtml: {
      "en": "<a target='_blank' href='https://developer.godaddy.com/keys'>Create API KEY</a>",
      "zh-cn": "<a target='_blank' href='https://developer.godaddy.com/keys'>创建 API KEY</a>",
    }
  },
  namecheap: {
    name: {
      "en": "Namecheap",
    },
    idLabel: "",
    secretLabel: "Password",
    helpHtml: {
      "en": "<a target='_blank' href='https://www.namecheap.com/support/knowledgebase/article.aspx/36/11/how-do-i-start-using-dynamic-dns/'>How to get started</a> <span style='color: red'>Namecheap DDNS does not support updating IPv6</span>",
      "zh-cn": "<a target='_blank' href='https://www.namecheap.com/support/knowledgebase/article.aspx/36/11/how-do-i-start-using-dynamic-dns/'>开启namecheap动态域名解析</a> <span style='color: red'>Namecheap DDNS 不支持更新 IPv6</span>",
    }
  },
  namesilo: {
    name: {
      "en": "NameSilo",
    },
    idLabel: "",
    secretLabel: "Password",
    helpHtml: {
      "en": "<a target='_blank' href='https://www.namesilo.com/account/api-manager'>How to get started</a> <b>Please note that the TTL of namesilo is at least 1 hour</b>",
      "zh-cn": "<a target='_blank' href='https://www.namesilo.com/account/api-manager'>开启namesilo动态域名解析</a> <b>请注意namesilo的TTL最低1小时</b>",
    }
  },
  vercel: {
    name: {
      "en": "Vercel",
    },
    idLabel: "",
    secretLabel: "Token",
    helpHtml: {
      "en": "<a target='_blank' href='https://vercel.com/account/tokens'>Create Token</a>",
      "zh-cn": "<a target='_blank' href='https://vercel.com/account/tokens'>创建令牌</a>",
    }
  },
  dynadot: {
    name: {
      "en": "Dynadot",
    },
    idLabel: "",
    secretLabel: "Password",
    helpHtml: {
      "en": "<a target='_blank' href='https://www.dynadot.com/community/help/question/enable-DDNS'>How to get started</a>",
      "zh-cn": "<a target='_blank' href='https://www.dynadot.com/community/help/question/enable-DDNS'>开启Dynadot动态域名解析</a>",
    }
  },
  trafficroute: {
    name: {
      "en": "TrafficRoute",
      "zh-cn": "火山引擎",
    },
    idLabel: "AccessKey",
    secretLabel: "SecretAccessKey",
    helpHtml: {
      "en": "<a target='_blank' href='https://console.volcengine.com/iam/keymanage/'>Create AccessKey</a>",
      "zh-cn": "<a target='_blank' href='https://console.volcengine.com/iam/keymanage/'>创建火山引擎 API 密钥</a>",
    }
  },
};

const SVG_CODE = {
  success: `<svg viewBox="64 64 896 896" focusable="false" data-icon="check-circle" width="1em" height="1em" fill="#52c41a" aria-hidden="true"><path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm193.5 301.7l-210.6 292a31.8 31.8 0 01-51.7 0L318.5 484.9c-3.8-5.3 0-12.7 6.5-12.7h46.9c10.2 0 19.9 4.9 25.9 13.3l71.2 98.8 157.2-218c6-8.3 15.6-13.3 25.9-13.3H699c6.5 0 10.3 7.4 6.5 12.7z"></path></svg>`,
  info: `<svg viewBox="64 64 896 896" focusable="false" data-icon="info-circle" width="1em" height="1em" fill="#1677ff" aria-hidden="true"><path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm32 664c0 4.4-3.6 8-8 8h-48c-4.4 0-8-3.6-8-8V456c0-4.4 3.6-8 8-8h48c4.4 0 8 3.6 8 8v272zm-32-344a48.01 48.01 0 010-96 48.01 48.01 0 010 96z"></path></svg>`,
  warning: '<svg viewBox="64 64 896 896" focusable="false" data-icon="exclamation-circle" width="1em" height="1em" fill="#faad14" aria-hidden="true"><path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm-32 232c0-4.4 3.6-8 8-8h48c4.4 0 8 3.6 8 8v272c0 4.4-3.6 8-8 8h-48c-4.4 0-8-3.6-8-8V296zm32 440a48.01 48.01 0 010-96 48.01 48.01 0 010 96z"></path></svg>',
  error: '<svg viewBox="64 64 896 896" focusable="false" data-icon="close-circle" width="1em" height="1em" fill="#ff4d4f" aria-hidden="true"><path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm165.4 618.2l-66-.3L512 563.4l-99.3 118.4-66.1.3c-4.4 0-8-3.5-8-8 0-1.9.7-3.7 1.9-5.2l130.1-155L340.5 359a8.32 8.32 0 01-1.9-5.2c0-4.4 3.6-8 8-8l66.1.3L512 464.6l99.3-118.4 66-.3c4.4 0 8 3.5 8 8 0 1.9-.7 3.7-1.9 5.2L553.5 514l130 155c1.2 1.5 1.9 3.3 1.9 5.2 0 4.4-3.6 8-8 8z"></path></svg>'
}


const I18N_MAP = {
  'en': {
    'Logs': 'Logs',
    'Save': 'Save',
    'Config:': 'Config:',
    'Add': 'Add',
    'Rename': 'Rename',
    'RenameHelp': 'Enter a new name:',
    'Delete': 'Delete',
    'DNS Provider': 'DNS Provider',
    'Create AccessKey': 'Create AccessKey',
    'Auto': 'Auto',
    '1s': '1s',
    '5s': '5s',
    '10s': '10s',
    '1m': '1m',
    '2m': '2m',
    '10m': '10m',
    '30m': '30m',
    '1h': '1h',
    'ttlHelp': 'You can modify it if the account supports a smaller TTL. The TTL will only be updated when the IP changes',
    'Enabled': 'Enabled',
    'Get IP method': 'Get IP method',
    'By api': 'By api',
    'By network card': 'By network card',
    'By command': 'By command',
    'domainsHelp': `
      Enter one domain per line.
      If the domain is unregistrable, manually separate it into a subdomain and a root domain by using a colon. e.g. <code>www:domain.example.com</code><br />

      Support for <a target="blank" href="https://github.com/jeessy2/ddns-go/wiki/传递自定义参数">custom parameters</a> (Simplified Chinese)
    `,
    'Regular exp.': 'Regular exp.',
    'regHelp': 'You can use @1 to specify the first IPv6 address, @2 to specify the second IPv6 address... You can also use regular expressions to match the specified IPv6 address, leave it blank to disable it',
    'Others': 'Others',
    'Deny from WAN': 'Deny from WAN',
    'NotAllowWanAccessHelp': 'Enable to deny access from the public network',
    'Username': 'Username',
    'accountHelp': 'Username/Password is required',
    'passwordHelp': 'If you need to change the password, please enter it here',
    'Password': 'Password',
    'WebhookURLHelp': `
      <a
        target="blank"
        href="https://github.com/jeessy2/ddns-go/blob/master/README_EN.md#webhook"
      >Click to get more info</a
      ><br />
      Support variables #{ipv4Addr}, #{ipv4Result},
      #{ipv4Domains}, #{ipv6Addr}, #{ipv6Result}, #{ipv6Domains}
    `,
    'WebhookRequestBodyHelp': 'If RequestBody is empty, it is a GET request, otherwise it is a POST request. Supported variables are the same as above',
    'WebhookHeadersHelp': 'One header per line, such as: Authorization: Bearer API_KEY',
    'Try it': 'Try it',
    'Clear': 'Clear',
    'OK': 'OK',
    "Ipv4UrlHelp": "https://api.ipify.org, https://myip.ipip.net, https://ddns.oray.com/checkip, https://ip.3322.net",
    "Ipv6UrlHelp": "https://speed.neu6.edu.cn/getIP.php, https://v6.ident.me, https://6.ipw.cn",
    "Ipv4NetInterfaceHelp": "Get IPv4 address through network card",
    "Ipv6NetInterfaceHelp": "If you do not specify a matching regular expression, the first IPv6 address will be used by default",
    "Ipv4CmdHelp": "Get IPv4 through command, only use the first matching IPv4 address of standard output(stdout). Such as: ip -4 addr show eth1",
    "Ipv6CmdHelp": "Get IPv6 through command, only use the first matching IPv6 address of standard output(stdout). Such as: ip -6 addr show eth1",
    "NetInterfaceEmptyHelp": '<span style="color: red">No available network card found</span>',
    "Login": 'Login',
    "LoginInit": 'Login and configure as an administrator account',
  },
  'zh-cn': {
    'Logs': '日志',
    'Save': '保存',
    'Config:': '配置切换:',
    'Add': '添加',
    'Rename': '重命名',
    'RenameHelp': '输入新名称：',
    'Delete': '删除',
    'DNS Provider': 'DNS服务商',
    'Create AccessKey': '创建 AccessKey',
    'Auto': '自动',
    '1s': '1秒',
    '5s': '5秒',
    '10s': '10秒',
    '1m': '1分钟',
    '2m': '2分钟',
    '10m': '10分钟',
    '30m': '30分钟',
    '1h': '1小时',
    'ttlHelp': '如账号支持更小的 TTL, 可修改。IP 有变化时才会更新TTL',
    'Enabled': '是否启用',
    'Get IP method': '获取 IP 方式',
    'By api': '通过接口获取',
    'By network card': '通过网卡获取',
    'By command': '通过命令获取',
    'domainsHelp': `
      每行一个域名。
      如果域名不可注册，请使用冒号手动将其分为子域名和根域名。如 <code>www:domain.example.com</code><br />

      支持<a target="blank" href="https://github.com/jeessy2/ddns-go/wiki/传递自定义参数">自定义参数</a>
    `,
    'Regular exp.': '匹配正则表达式',
    'regHelp': '可使用 @1 指定第一个IPv6地址, @2 指定第二个IPv6地址... 也可使用正则表达式匹配指定的IPv6地址, 留空则不启用',
    'Others': '其他',
    'Deny from WAN': '禁止公网访问',
    'NotAllowWanAccessHelp': '启用后禁止从公网访问此页面',
    'Username': '用户名',
    'accountHelp': '必须输入用户名/密码',
    'passwordHelp': '如需修改密码，请在此处输入新密码',
    'Password': '密码',
    'WebhookURLHelp': `
      <a target="blank" href="https://github.com/jeessy2/ddns-go#webhook">点击参考官方 Webhook 说明</a>
      <br />
      支持的变量 #{ipv4Addr}, #{ipv4Result}, #{ipv4Domains}, #{ipv6Addr}, #{ipv6Result}, #{ipv6Domains}
    `,
    'WebhookRequestBodyHelp': '如果 RequestBody 为空, 则为 GET 请求, 否则为 POST 请求。支持的变量同上',
    'WebhookHeadersHelp': '一行一个Header, 如: Authorization: Bearer API_KEY',
    'Try it': '模拟测试Webhook',
    'Clear': '清空',
    'OK': '确定',
    "Ipv4UrlHelp": "https://myip.ipip.net, https://ddns.oray.com/checkip, https://ip.3322.net",
    "Ipv6UrlHelp": "https://speed.neu6.edu.cn/getIP.php, https://v6.ident.me, https://6.ipw.cn",
    "Ipv4NetInterfaceHelp": "通过网卡获取IPv4",
    "Ipv6NetInterfaceHelp": "如不指定匹配正则表达式，将默认使用第一个 IPv6 地址",
    "Ipv4CmdHelp": `
      通过命令获取IPv4, 仅使用标准输出(stdout)的第一个匹配的 IPv4 地址。如: ip -4 addr show eth1
      <a target="blank" href="https://github.com/jeessy2/ddns-go/wiki/通过命令获取IP参考">点击参考更多</a>
    `,
    "Ipv6CmdHelp": `
      通过命令获取IPv6, 仅使用标准输出(stdout)的第一个匹配的 IPv6 地址。如: ip -6 addr show eth1
      <a target="blank" href="https://github.com/jeessy2/ddns-go/wiki/通过命令获取IP参考">点击参考更多</a>
    `,
    "NetInterfaceEmptyHelp": '<span style="color: red">没有找到可用的网卡</span>',
    "Login": '登录',
    "LoginInit": '登录并配置为管理员账号',
  }
};
