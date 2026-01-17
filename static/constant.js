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
  aliesa: {
    name: {
      "en": "Aliyun ESA",
      "zh-cn": "阿里云 ESA",
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
      "en": "<a target='_blank' href='https://console.bce.baidu.com/iam/?_=1651763238057#/iam/accesslist'>Create AccessKey</a>",
      "zh-cn": "<a target='_blank' href='https://console.bce.baidu.com/iam/?_=1651763238057#/iam/accesslist'>创建 AccessKey</a>",
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
      "en": "<a target='_blank' href='https://developer.godaddy.com/keys'>Create API KEY</a><br/><span style='color: #ff9800;'>⚠️ Note: GoDaddy API requires you to have 10 or more domains or a Pro plan</span>",
      "zh-cn": "<a target='_blank' href='https://developer.godaddy.com/keys'>创建 API KEY</a><br/><span style='color: #ff9800;'>⚠️ 温馨提示：GoDaddy 现在需要拥有 10 个及以上的域名或 Pro Plan 才可以使用 API</span>",
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
    },
    extParamLabel: "Team ID",
    extParamHelpHtml: {
      "en": "Optional. If you are using a Vercel Team account, please fill in the Team ID",
      "zh-cn": "可选项，如果您使用的是 Vercel 团队账户，请填写团队 ID"
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
  dynv6: {
    name: {
      "en": "Dynv6",
    },
    idLabel: "",
    secretLabel: "Token",
    helpHtml: {
      "en": "<a target='_blank' href='https://dynv6.com/keys'>Create Token</a>",
      "zh-cn": "<a target='_blank' href='https://dynv6.com/keys'>创建令牌</a>",
    }
  },
  spaceship: {
    name: {
      "en": "Spaceship",
    },
    idLabel: "API Key",
    secretLabel: "API Secret",
    helpHtml: {
      "en": "<a target='_blank' href='https://www.spaceship.com/application/api-manager/'>Create API Key</a>",
      "zh-cn": "<a target='_blank' href='https://www.spaceship.com/application/api-manager/'>创建 API 密钥</a>",
    }
  },
  dnsla: {
    name: {
      "en": "Dnsla",
      "zh-cn": "Dnsla",
    },
    idLabel: "APIID",
    secretLabel: "API密钥",
    helpHtml: {
      "en": "<a target='_blank' href='https://console.dns.la/login?aksk=1'>Create AccessKey</a>",
      "zh-cn": "<a target='_blank' href='https://console.dns.la/login?aksk=1'>创建 AccessKey</a>",
    }
  },
  nowcn: {
    name: {
      "en": "Nowcn",
      "zh-cn": "时代互联",
    },
    idLabel: "auth-userid",
    secretLabel: "api-key",
    helpHtml: {
      "en": "<a target='_blank' href='https://www.now.cn/'>api-key</a>",
      "zh-cn": "<a target='_blank' href='https://www.now.cn/'>获取 api-key</a>",
    }
  },
  eranet: {
    name: {
      "en": "Eranet",
      "zh-cn": "Eranet",
    },
    idLabel: "auth-userid",
    secretLabel: "api-key",
    helpHtml: {
      "en": "<a target='_blank' href='https://partner.eranet.com/admin/mode_Http_Api_detail.php'>api-key</a>",
      "zh-cn": "<a target='_blank' href='https://partner.eranet.com/admin/mode_Http_Api_detail.php'>获取 api-key</a>",
    }
  },
  gcore: {
    name: {
      "en": "Gcore",
    },
    idLabel: "",
    secretLabel: "API Token",
    helpHtml: {
      "en": "<a target='_blank' href='https://portal.gcore.com/accounts/profile/api-tokens/create'>Create API Token</a>",
      "zh-cn": "<a target='_blank' href='https://portal.gcore.com/accounts/profile/api-tokens/create'>创建 API Token</a>",
    }
  },
  edgeone: {
    name: {
      "en": "Edgeone",
      "zh-cn": "Edgeone",
    },
    idLabel: "SecretId",
    secretLabel: "SecretKey",
    helpHtml: {
      "en": "<a target='_blank' href='https://console.cloud.tencent.com/cam/capi'>Create AccessKey</a>",
      "zh-cn": "<a target='_blank' href='https://console.cloud.tencent.com/cam/capi'>创建腾讯云 API 密钥</a>",
    }
  },
  nsone: {
    name: {
      "en": "IBM NS1 Connect",
      "zh-cn": "IBM NS1 Connect",
    },
    idLabel: "",
    secretLabel: "API Key",
    helpHtml: {
      "en": "<a target='_blank' href='https://my.nsone.net/#/account/settings/keys'>Create API Key</a>",
      "zh-cn": "<a target='_blank' href='https://my.nsone.net/#/account/settings/keys'>创建 API 密钥</a>",
    }
  },
  name_com: {
    name: {
      "en": "name.com",
      "zh-cn": "name.com",
    },
    idLabel: "username",
    secretLabel: "token",
    helpHtml: {
      "en": "<a target='_blank' href='https://www.name.com/zh-cn/account/settings/api'>name.com Create API Token</a>",
      "zh-cn": "<a target='_blank' href='https://www.name.com/zh-cn/account/settings/api'>name.com 创建 API Token</a>",
    }
  },
};

const SVG_CODE = {
  success: `<svg viewBox="64 64 896 896" focusable="false" data-icon="check-circle" width="1em" height="1em" fill="#52c41a" aria-hidden="true"><path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm193.5 301.7l-210.6 292a31.8 31.8 0 01-51.7 0L318.5 484.9c-3.8-5.3 0-12.7 6.5-12.7h46.9c10.2 0 19.9 4.9 25.9 13.3l71.2 98.8 157.2-218c6-8.3 15.6-13.3 25.9-13.3H699c6.5 0 10.3 7.4 6.5 12.7z"></path></svg>`,
  info: `<svg viewBox="64 64 896 896" focusable="false" data-icon="info-circle" width="1em" height="1em" fill="#1677ff" aria-hidden="true"><path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm32 664c0 4.4-3.6 8-8 8h-48c-4.4 0-8-3.6-8-8V456c0-4.4 3.6-8 8-8h48c4.4 0 8 3.6 8 8v272zm-32-344a48.01 48.01 0 010-96 48.01 48.01 0 010 96z"></path></svg>`,
  warning: '<svg viewBox="64 64 896 896" focusable="false" data-icon="exclamation-circle" width="1em" height="1em" fill="#faad14" aria-hidden="true"><path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm-32 232c0-4.4 3.6-8 8-8h48c4.4 0 8 3.6 8 8v272c0 4.4-3.6 8-8 8h-48c-4.4 0-8-3.6-8-8V296zm32 440a48.01 48.01 0 010-96 48.01 48.01 0 010 96z"></path></svg>',
  error: '<svg viewBox="64 64 896 896" focusable="false" data-icon="close-circle" width="1em" height="1em" fill="#ff4d4f" aria-hidden="true"><path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm165.4 618.2l-66-.3L512 563.4l-99.3 118.4-66.1.3c-4.4 0-8-3.5-8-8 0-1.9.7-3.7 1.9-5.2l130.1-155L340.5 359a8.32 8.32 0 01-1.9-5.2c0-4.4 3.6-8 8-8l66.1.3L512 464.6l99.3-118.4 66-.3c4.4 0 8 3.5 8 8 0 1.9-.7 3.7-1.9 5.2L553.5 514l130 155c1.2 1.5 1.9 3.3 1.9 5.2 0 4.4-3.6 8-8 8z"></path></svg>'
}
