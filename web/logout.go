package web

import (
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	// 创建一个过期的 Cookie 来清除客户端的身份认证 Cookie
	expiredCookie := http.Cookie{
		Name:     cookieName, // 假设你的身份验证使用的是名为 "auth" 的 Cookie
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // 设置为过期时间
		MaxAge:   -1,              // 立即删除该 Cookie
		HttpOnly: true,
	}
	// 设置过期的 Cookie
	http.SetCookie(w, &expiredCookie)

	// 重定向用户到登录页面
	http.Redirect(w, r, "/login", http.StatusFound)
}
