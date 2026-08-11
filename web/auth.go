package web

import (
	"net/http"
	"os"
	"time"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

// TrustedAuthHeaderEnv, when set, names a request header injected by a trusted
// reverse proxy performing forward-auth (e.g. Authentik, Authelia, Traefik).
// If that header is present on a request coming from a private network, the
// login check is bypassed so the proxy's authentication is trusted. Empty or
// unset (the default) keeps the original behavior unchanged. The private-network
// guard prevents public clients from spoofing the header.
const TrustedAuthHeaderEnv = "DDNS_GO_TRUSTED_AUTH_HEADER"

// ViewFunc func
type ViewFunc func(http.ResponseWriter, *http.Request)

// Auth 验证Token是否已经通过
func Auth(f ViewFunc) ViewFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 信任反向代理注入的认证头(forward-auth), 仅对来自私网的请求生效, 防止公网伪造
		if h := os.Getenv(TrustedAuthHeaderEnv); h != "" && r.Header.Get(h) != "" && util.IsPrivateNetwork(r.RemoteAddr) {
			f(w, r)
			return
		}
		cookieInWeb, err := r.Cookie(cookieName)
		if err != nil {
			http.Redirect(w, r, "./login", http.StatusTemporaryRedirect)
			return
		}

		conf, _ := config.GetConfigCached()

		// 禁止公网访问
		if conf.NotAllowWanAccess {
			if !util.IsPrivateNetwork(r.RemoteAddr) {
				w.WriteHeader(http.StatusForbidden)
				util.Log("%q 被禁止从公网访问", util.GetRequestIPStr(r))
				return
			}
		}

		// 验证token
		if cookieInSystem.Value != "" &&
			cookieInSystem.Value == cookieInWeb.Value &&
			cookieInSystem.Expires.After(time.Now()) {
			f(w, r) // 执行被装饰的函数
			return
		}

		http.Redirect(w, r, "./login", http.StatusTemporaryRedirect)
	}
}

// AuthAssert 保护静态等文件不被公网访问
func AuthAssert(f ViewFunc) ViewFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		conf, err := config.GetConfigCached()

		// 配置文件为空, 启动时间超过3小时禁止从公网访问
		if err != nil &&
			time.Since(startTime) > time.Duration(3*time.Hour) && !util.IsPrivateNetwork(r.RemoteAddr) {
			w.WriteHeader(http.StatusForbidden)
			util.Log("%q 配置文件为空, 超过3小时禁止从公网访问", util.GetRequestIPStr(r))
			return
		}

		// 禁止公网访问
		if conf.NotAllowWanAccess {
			if !util.IsPrivateNetwork(r.RemoteAddr) {
				w.WriteHeader(http.StatusForbidden)
				util.Log("%q 被禁止从公网访问", util.GetRequestIPStr(r))
				return
			}
		}

		f(w, r) // 执行被装饰的函数

	}
}
