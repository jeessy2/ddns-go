function updateColorSchemeMeta(isDark) {
  const meta = document.querySelector('meta[name="color-scheme"]');
  if (meta) {
    meta.setAttribute('content', isDark ? 'dark' : 'light');
  }
}

function toggleTheme(write = false) {
  const docEle = document.documentElement;
  if (docEle.getAttribute("data-theme") === "dark") {
    docEle.removeAttribute("data-theme");
    updateColorSchemeMeta(false);
    write && localStorage.setItem("theme", "light");
  } else {
    docEle.setAttribute("data-theme", "dark");
    updateColorSchemeMeta(true);
    write && localStorage.setItem("theme", "dark");
  }
}

const theme = localStorage.getItem("theme") ??
  (window.matchMedia("(prefers-color-scheme: dark)").matches
    ? "dark"
    : "light");

if (theme === "dark") {
  toggleTheme();
}

// 长按重置功能
let pressTimer = null;
let isLongPress = false;

const button = document.getElementById("themeButton");

function startPress() {
  isLongPress = false;
  // 800ms后触发长按
  pressTimer = setTimeout(() => {
    isLongPress = true;

    // 清除用户偏好，恢复自动模式
    localStorage.removeItem("theme");

    // 立即同步系统主题状态
    const systemIsDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
    const currentlyDark = document.documentElement.getAttribute("data-theme") === "dark";
    if (systemIsDark !== currentlyDark) {
      toggleTheme();
    }

    // 显示成功提示
    showMessage({
      content: i18n({
        "en": "Theme has been restored to auto mode",
        "zh-cn": "主题已恢复自动跟随系统"
      }),
      type: "success",
      duration: 2000
    });
  }, 800);
}

function endPress() {
  clearTimeout(pressTimer);
  // 短按才执行切换
  if (!isLongPress) {
    toggleTheme(true);
  }
}

function cancelPress() {
  clearTimeout(pressTimer);
}

// 鼠标事件
button.addEventListener('mousedown', startPress);
button.addEventListener('mouseup', endPress);
button.addEventListener('mouseleave', cancelPress);

// 触摸事件（移动设备）
button.addEventListener('touchstart', (e) => {
  e.preventDefault(); // 防止触发点击
  startPress();
});
button.addEventListener('touchmove', cancelPress);
button.addEventListener('touchend', endPress);
button.addEventListener('touchcancel', cancelPress);

// 系统主题变化监听器
// 仅在自动模式下响应（即用户未手动设置偏好时）
window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", (e) => {
  if (!localStorage.getItem("theme")) {
    // 只有在没有用户偏好时才自动切换
    const shouldBeDark = e.matches;
    const currentlyDark = document.documentElement.getAttribute("data-theme") === "dark";
    if (shouldBeDark !== currentlyDark) {
      toggleTheme();
    }
  }
});
