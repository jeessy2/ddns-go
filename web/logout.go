package web

import (
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	// 覆盖cookieInSystem
	cookieInSystem = &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // 设置为过期时间
		MaxAge:   -1,              // 立即删除该 Cookie
		HttpOnly: true,
	}
	// 设置过期的 Cookie
	http.SetCookie(w, cookieInSystem)

	// 重定向用户到登录页面
	http.Redirect(w, r, "./login", http.StatusFound)
}
