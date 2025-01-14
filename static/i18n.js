const I18N_MAP = {
  'Logs': {
    'en': 'Logs',
    'zh-cn': '日志'
  },
  'Save': {
    'en': 'Save',
    'zh-cn': '保存'
  },
  'Config:': {
    'en': 'Config:',
    'zh-cn': '配置切换:'
  },
  'Add': {
    'en': 'Add',
    'zh-cn': '添加'
  },
  'Rename': {
    'en': 'Rename',
    'zh-cn': '重命名'
  },
  'RenameHelp': {
    'en': 'Enter a new name:',
    'zh-cn': '输入新名称：'
  },
  'Delete': {
    'en': 'Delete',
    'zh-cn': '删除'
  },
  'DNS Provider': {
    'en': 'DNS Provider',
    'zh-cn': 'DNS服务商'
  },
  'Create AccessKey': {
    'en': 'Create AccessKey',
    'zh-cn': '创建 AccessKey'
  },
  'Auto': {
    'en': 'Auto',
    'zh-cn': '自动'
  },
  '1s': {
    'en': '1s',
    'zh-cn': '1秒'
  },
  '5s': {
    'en': '5s',
    'zh-cn': '5秒'
  },
  '10s': {
    'en': '10s',
    'zh-cn': '10秒'
  },
  '1m': {
    'en': '1m',
    'zh-cn': '1分钟'
  },
  '2m': {
    'en': '2m',
    'zh-cn': '2分钟'
  },
  '10m': {
    'en': '10m',
    'zh-cn': '10分钟'
  },
  '30m': {
    'en': '30m',
    'zh-cn': '30分钟'
  },
  '1h': {
    'en': '1h',
    'zh-cn': '1小时'
  },
  'ttlHelp': {
    'en': 'You can modify it if the account supports a smaller TTL. The TTL will only be updated when the IP changes',
    'zh-cn': '如账号支持更小的 TTL, 可修改。IP 有变化时才会更新TTL'
  },
  'Enabled': {
    'en': 'Enabled',
    'zh-cn': '是否启用'
  },
  'Get IP method': {
    'en': 'Get IP method',
    'zh-cn': '获取 IP 方式'
  },
  'By api': {
    'en': 'By api',
    'zh-cn': '通过接口获取'
  },
  'By network card': {
    'en': 'By network card',
    'zh-cn': '通过网卡获取'
  },
  'By command': {
    'en': 'By command',
    'zh-cn': '通过命令获取'
  },
  'domainsHelp': {
    'en': `
      Enter one domain per line.
      If the domain is unregistrable, manually separate it into a subdomain and a root domain by using a colon. e.g. <code>www:domain.example.com</code><br />

      Support for <a target="blank" href="https://github.com/jeessy2/ddns-go/wiki/传递自定义参数">custom parameters</a> (Simplified Chinese)
    `,
    'zh-cn': `
      每行一个域名。
      如果域名不可注册，请使用冒号手动将其分为子域名和根域名。如 <code>www:domain.example.com</code><br />
      支持<a target="blank" href="https://github.com/jeessy2/ddns-go/wiki/传递自定义参数">自定义参数</a>
    `
  },
  'Regular exp.': {
    'en': 'Regular exp.',
    'zh-cn': '匹配正则表达式'
  },
  'regHelp': {
    'en': 'You can use @1 to specify the first IPv6 address, @2 to specify the second IPv6 address... You can also use regular expressions to match the specified IPv6 address, leave it blank to disable it',
    'zh-cn': '可使用 @1 指定第一个IPv6地址, @2 指定第二个IPv6地址... 也可使用正则表达式匹配指定的IPv6地址, 留空则不启用'
  },
  'Others': {
    'en': 'Others',
    'zh-cn': '其他'
  },
  'Deny from WAN': {
    'en': 'Deny from WAN',
    'zh-cn': '禁止公网访问'
  },
  'NotAllowWanAccessHelp': {
    'en': 'Enable to deny access from the public network',
    'zh-cn': '启用后禁止从公网访问此页面'
  },
  'Username': {
    'en': 'Username',
    'zh-cn': '用户名'
  },
  'accountHelp': {
    'en': 'Username/Password is required',
    'zh-cn': '必须输入用户名/密码'
  },
  'passwordHelp': {
    'en': 'If you need to change the password, please enter it here',
    'zh-cn': '如需修改密码，请在此处输入新密码'
  },
  'Password': {
    'en': 'Password',
    'zh-cn': '密码'
  },
  'WebhookURLHelp': {
    'en': `
      <a
        target="blank"
        href="https://github.com/jeessy2/ddns-go/blob/master/README_EN.md#webhook"
      >Click to get more info</a
      ><br />
      Support variables #{ipv4Addr}, #{ipv4Result},
      #{ipv4Domains}, #{ipv6Addr}, #{ipv6Result}, #{ipv6Domains}
    `,
    'zh-cn': `
      <a target="blank" href="https://github.com/jeessy2/ddns-go#webhook">点击参考官方 Webhook 说明</a>
      <br />
      支持的变量 #{ipv4Addr}, #{ipv4Result}, #{ipv4Domains}, #{ipv6Addr}, #{ipv6Result}, #{ipv6Domains}
    `
  },
  'WebhookRequestBodyHelp': {
    'en': 'If RequestBody is empty, it is a GET request, otherwise it is a POST request. Supported variables are the same as above',
    'zh-cn': '如果 RequestBody 为空, 则为 GET 请求, 否则为 POST 请求。支持的变量同上'
  },
  'WebhookHeadersHelp': {
    'en': 'One header per line, such as: Authorization: Bearer API_KEY',
    'zh-cn': '一行一个Header, 如: Authorization: Bearer API_KEY'
  },
  'Try it': {
    'en': 'Try it',
    'zh-cn': '模拟测试Webhook'
  },
  'Clear': {
    'en': 'Clear',
    'zh-cn': '清空'
  },
  'OK': {
    'en': 'OK',
    'zh-cn': '确定'
  },
  "Ipv4UrlHelp": {
    'en': "https://api.ipify.org, https://myip.ipip.net, https://ddns.oray.com/checkip, https://ip.3322.net, https://v4.yinghualuo.cn/bejson",
    'zh-cn': "https://myip.ipip.net, https://ddns.oray.com/checkip, https://ip.3322.net, https://v4.yinghualuo.cn/bejson"
  },
  "Ipv6UrlHelp": {
    'en': "https://speed.neu6.edu.cn/getIP.php, https://v6.ident.me, https://6.ipw.cn, https://v6.yinghualuo.cn/bejson",
    'zh-cn': "https://speed.neu6.edu.cn/getIP.php, https://v6.ident.me, https://6.ipw.cn, https://v6.yinghualuo.cn/bejson"
  },
  "Ipv4NetInterfaceHelp": {
    'en': "Get IPv4 address through network card",
    'zh-cn': "通过网卡获取IPv4"
  },
  "Ipv6NetInterfaceHelp": {
    'en': "If you do not specify a matching regular expression, the first IPv6 address will be used by default",
    'zh-cn': "如不指定匹配正则表达式，将默认使用第一个 IPv6 地址"
  },
  "Ipv4CmdHelp": {
    'en': "Get IPv4 through command, only use the first matching IPv4 address of standard output(stdout). Such as: ip -4 addr show eth1",
    'zh-cn': `
      通过命令获取IPv4, 仅使用标准输出(stdout)的第一个匹配的 IPv4 地址。如: ip -4 addr show eth1
      <a target="blank" href="https://github.com/jeessy2/ddns-go/wiki/通过命令获取IP参考">点击参考更多</a>
    `
  },
  "Ipv6CmdHelp": {
    'en': "Get IPv6 through command, only use the first matching IPv6 address of standard output(stdout). Such as: ip -6 addr show eth1",
    'zh-cn': `
      通过命令获取IPv6, 仅使用标准输出(stdout)的第一个匹配的 IPv6 地址。如: ip -6 addr show eth1
      <a target="blank" href="https://github.com/jeessy2/ddns-go/wiki/通过命令获取IP参考">点击参考更多</a>
    `
  },
  "NetInterfaceEmptyHelp": {
    'en': '<span style="color: red">No available network card found</span>',
    'zh-cn': '<span style="color: red">没有找到可用的网卡</span>'
  },
  "Login": {
    'en': 'Login',
    'zh-cn': '登录'
  },
  "LoginInit": {
    'en': 'Login and configure as an administrator account',
    'zh-cn': '登录并配置为管理员账号'
  },
  "Logout": {
    'en': 'Logout',
    'zh-cn': '注销'
  },
  "webhookTestTooltip": {
    'en': 'Send a fake data to the Webhook URL immediately to test if the Webhook is working properly',
    'zh-cn': '立即发送一条假数据到Webhook URL，用于测试Webhook是否正常工作'
  },
  "themeTooltip": {
    'en': 'Switch between light and dark themes',
    'zh-cn': '切换明暗主题'
  },
};

const LANG = localStorage.getItem('lang') || (navigator.language || navigator.browserLanguage).replaceAll('_', '-').toLowerCase();

const getLocalLang = (langs) => {
  // 优先取地区语言
  if (langs.includes(LANG)) {
    return LANG;
  }
  // 其次取表示语言
  if (langs.includes(LANG.split('-')[0])) {
    return LANG.split('-')[0];
  }
  // 再取表示语言相同的地区语言
  for (const l of langs) {
    if (l.split('-')[0] === LANG.split('-')[0]) {
      return l;
    }
  }
  // 无法匹配则取英文
  return 'en';
}

// 支持两种调用方式：
// 1. 文本在I18N字典中的key，如"hello"
// 2. 语言字符串字典，{en: "hello", zh: "你好"}
const i18n = (keyOrLangDict) => {
  let key = keyOrLangDict;
  let langDict = keyOrLangDict;
  if (typeof keyOrLangDict === 'string') {
    langDict = I18N_MAP[keyOrLangDict];
  } else {
    key = null;
  }
  if (!langDict) {
    console.warn(`i18n: No translation for key "${key}"`);
    return key;
  }
  const lang = getLocalLang(Object.keys(langDict));
  if (lang in langDict) {
    return langDict[lang];
  }
  console.warn(`i18n: No such language "${lang}" in langDict ${langDict}`);
  return key;
}

const convertDom = (dom = document) => {
  dom.querySelectorAll('[data-i18n]').forEach(el => {
    const key = el.dataset.i18n;
    el.textContent = i18n(key);
  });
  dom.querySelectorAll('[data-i18n-html]').forEach(el => {
    const key = el.dataset.i18nHtml;
    el.innerHTML = i18n(key);
  });
  dom.querySelectorAll('[data-i18n-attr]').forEach(el => {
    el.dataset.i18nAttr.split(',').forEach(item => {
      let [attr, key] = item.split(':');
      attr = attr.trim();
      key = key || el.getAttribute(attr);
      el.setAttribute(attr, i18n(key));
    });
  });
}

document.addEventListener('DOMContentLoaded', () => {convertDom();});