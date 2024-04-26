const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

const html2Element = (htmlString) => {
  const doc = new DOMParser().parseFromString(htmlString, 'text/html')
  return doc.body.firstElementChild
}

// 在页面顶部显示一行消息，并在若干秒后自动消失
const showMessage = async (msgObj) => {
  // 填充默认值
  msgObj = Object.assign({
    type: 'info',
    content: '',
    html: false,
    duration: 3000
  }, msgObj)
  // 当前是否有消息容器
  let $container = document.getElementById('msg-container')
  if (!$container) {
    // 创建消息容器
    $container = html2Element('<div id="msg-container"></div>')
    document.body.appendChild($container)
  }
  // 创建消息元素
  const $msg = html2Element('<div class="msg msg-fade"></div>')
  // 创建两个span，用于显示消息的图标和内容
  const $content = html2Element('<span></span>')

  // 填充内容，根据html属性决定使用text还是html
  if (msgObj.html) {
    $content.innerHTML = msgObj.content
  } else {
    $content.textContent = msgObj.content
  }
  // 根据消息类型设置图标
  $msg.innerHTML = `<span class="msg-icon">${SVG_CODE[msgObj.type]}</span>`
  $msg.appendChild($content)
  $container.appendChild($msg)
  // 确保动画生效
  await delay(0)
  $msg.classList.remove('msg-fade')
  // 等待动画结束
  await delay(200)
  // 销毁函数
  const destroy = async () => {
    // 增加消失动画
    $msg.classList.add('msg-fade')
    // 动画结束后移除元素
    await delay(200)
    $msg.remove()
    // 如果容器中没有消息了，移除容器
    if (!$container.children.length) {
      $container.remove()
    }
  }
  // 如果duration为0，则不自动消失
  if (msgObj.duration === 0) {
    return destroy
  }
  // 自动消失计时器
  let timer = setTimeout(destroy, msgObj.duration)
  // 注册鼠标事件，鼠标移入时取消自动消失
  $msg.addEventListener('mouseenter', () => {
    clearTimeout(timer)
  })
  // 鼠标移出时重新计时
  $msg.addEventListener('mouseleave', () => {
    timer = setTimeout(destroy, msgObj.duration)
  })
  return destroy
}

const request = {
  baseURL: './',
  parse: async function(resp) {
    const text = await resp.text()
    try {
      return JSON.parse(text)
    } catch (e) {
      return text
    }
  },
  stringify: function(dict) {
    const result = []
    for (let key in dict) {
      if (!dict.hasOwnProperty(key)) {
        continue
      }
      // 所有空值将被删除
      if (String(dict[key])) {
        result.push(`${key}=${encodeURIComponent(dict[key])}`)
      }
    }
    return result.join('&')
  },
  get: async function(path, data, parseFunc) {
    const response = await fetch(`${this.baseURL}${path}?${this.stringify(data)}`)
    if (response.redirected) {
      window.location.href = response.url
    }
    return await (parseFunc||this.parse)(response)
  },
  post: async function(path, data, parseFunc) {
    if (typeof data === 'object') {
      data = JSON.stringify(data)
    }
    const response = await fetch(`${this.baseURL}${path}`, {
      method: 'POST',
      body: data
    })
    if (response.redirected) {
      window.location.href = response.url
    }
    return await (parseFunc||this.parse)(response)
  }
}