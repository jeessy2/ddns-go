// 常量资源
const SVG_CODE = {
  success: `<svg viewBox="64 64 896 896" focusable="false" data-icon="check-circle" width="1em" height="1em" fill="#52c41a" aria-hidden="true"><path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm193.5 301.7l-210.6 292a31.8 31.8 0 01-51.7 0L318.5 484.9c-3.8-5.3 0-12.7 6.5-12.7h46.9c10.2 0 19.9 4.9 25.9 13.3l71.2 98.8 157.2-218c6-8.3 15.6-13.3 25.9-13.3H699c6.5 0 10.3 7.4 6.5 12.7z"></path></svg>`,
  info: `<svg viewBox="64 64 896 896" focusable="false" data-icon="info-circle" width="1em" height="1em" fill="#1677ff" aria-hidden="true"><path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm32 664c0 4.4-3.6 8-8 8h-48c-4.4 0-8-3.6-8-8V456c0-4.4 3.6-8 8-8h48c4.4 0 8 3.6 8 8v272zm-32-344a48.01 48.01 0 010-96 48.01 48.01 0 010 96z"></path></svg>`,
  warning: '<svg viewBox="64 64 896 896" focusable="false" data-icon="exclamation-circle" width="1em" height="1em" fill="#faad14" aria-hidden="true"><path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm-32 232c0-4.4 3.6-8 8-8h48c4.4 0 8 3.6 8 8v272c0 4.4-3.6 8-8 8h-48c-4.4 0-8-3.6-8-8V296zm32 440a48.01 48.01 0 010-96 48.01 48.01 0 010 96z"></path></svg>',
  error: '<svg viewBox="64 64 896 896" focusable="false" data-icon="close-circle" width="1em" height="1em" fill="#ff4d4f" aria-hidden="true"><path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm165.4 618.2l-66-.3L512 563.4l-99.3 118.4-66.1.3c-4.4 0-8-3.5-8-8 0-1.9.7-3.7 1.9-5.2l130.1-155L340.5 359a8.32 8.32 0 01-1.9-5.2c0-4.4 3.6-8 8-8l66.1.3L512 464.6l99.3-118.4 66-.3c4.4 0 8 3.5 8 8 0 1.9-.7 3.7-1.9 5.2L553.5 514l130 155c1.2 1.5 1.9 3.3 1.9 5.2 0 4.4-3.6 8-8 8z"></path></svg>'
}



const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

// 在页面顶部显示一行消息，并在若干秒后自动消失
const showMessage = async (msgObj) => {
  // 当前是否有消息容器
  let $container = document.querySelector('#msg-container')
  if (!$container) {
    // 创建消息容器
    $container = document.createElement('div')
    $container.id = 'msg-container'
    document.body.appendChild($container)
  }
  // 创建消息元素
  const $msg = document.createElement('div')
  // 创建两个span，用于显示消息的图标和内容
  const $icon = document.createElement('span')
  const $content = document.createElement('span')
  $icon.classList.add('msg-icon')
  // 根据消息类型设置图标
  $icon.innerHTML = SVG_CODE[msgObj.type] || SVG_CODE.info
  $content.innerText = msgObj.content || ''
  $msg.appendChild($icon)
  $msg.appendChild($content)
  // 增加出现动画
  $msg.classList.add('msg','msg-fade')
  $container.appendChild($msg)
  // 0延迟是为了让剩余的代码存入异步队列，稍后执行。否则浏览器会把两步操作合并，导致动画不生效
  await delay(0)
  $msg.classList.remove('msg-fade')
  // 等待动画结束
  await delay(200)
  // 消失函数
  const disappear = async () => {
    // 清除计时器
    clearTimeout(timer)
    // 增加消失动画
    $msg.classList.add('msg-fade')
    // 动画结束后移除元素
    await delay(200);
    $container.removeChild($msg)
    // 如果容器中没有消息了，移除容器
    if ($container.children.length === 0) {
      document.body.removeChild($container)
    }
  }
  // 如果duration为0，则不自动消失
  if (msgObj.duration === 0) {
    return disappear
  }
  // 自动消失计时器
  let timer = setTimeout(disappear, msgObj.duration || 3000)
  // 注册鼠标事件，鼠标移入时取消自动消失
  $msg.onmouseenter = () => {
    clearTimeout(timer)
  }
  // 鼠标移出时重新计时
  $msg.onmouseleave = () => {
    timer = setTimeout(disappear, msgObj.duration || 3000)
  }
  return disappear
}