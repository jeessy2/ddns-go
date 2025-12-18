function toggleTheme(write = false) {
  const docEle = document.documentElement;
  if (docEle.getAttribute("data-theme") === "dark") {
    docEle.removeAttribute("data-theme");
    write && localStorage.setItem("theme", "light");
  } else {
    docEle.setAttribute("data-theme", "dark");
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

// 主题切换
document.getElementById("themeButton").addEventListener('click', () => toggleTheme(true));
