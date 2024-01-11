package web

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

// ViewFunc func
type ViewFunc func(http.ResponseWriter, *http.Request)

type loginDetect struct {
	FailTimes int
}

var ld = &loginDetect{}

// BasicAuth basic auth
func BasicAuth(f ViewFunc) ViewFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conf, err := config.GetConfigCached()

		// 配置文件为空, 启动时间超过3小时禁止从公网访问
		if err != nil && time.Now().Unix()-startTime > 3*60*60 &&
			(!util.IsPrivateNetwork(r.RemoteAddr) || !util.IsPrivateNetwork(r.Host)) {
			w.WriteHeader(http.StatusForbidden)
			util.Log("%q 配置文件为空, 超过3小时禁止从公网访问", util.GetRequestIPStr(r))
			return
		}

		// 禁止公网访问
		if conf.NotAllowWanAccess {
			if !util.IsPrivateNetwork(r.RemoteAddr) || !util.IsPrivateNetwork(r.Host) {
				w.WriteHeader(http.StatusForbidden)
				util.Log("%q 被禁止从公网访问", util.GetRequestIPStr(r))
				return
			}
		}

		// 帐号或密码为空。跳过
		if conf.Username == "" && conf.Password == "" {
			// 执行被装饰的函数
			f(w, r)
			return
		}

		if ld.FailTimes >= 5 {
			util.Log("%q 登陆失败超过5次! 并延时5分钟响应", util.GetRequestIPStr(r))
			time.Sleep(5 * time.Minute)
			if ld.FailTimes >= 5 {
				ld.FailTimes = 0
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// 认证帐号密码
		basicAuthPrefix := "Basic "

		// 获取 request header
		auth := r.Header.Get("Authorization")
		// 如果是 http basic auth
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// 解码认证信息
			payload, err := base64.StdEncoding.DecodeString(
				auth[len(basicAuthPrefix):],
			)
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 &&
					bytes.Equal(pair[0], []byte(conf.Username)) &&
					bytes.Equal(pair[1], []byte(conf.Password)) {
					ld.FailTimes = 0
					// 执行被装饰的函数
					f(w, r)
					return
				}
			}

			ld.FailTimes = ld.FailTimes + 1
			util.Log("%q 帐号密码不正确", util.GetRequestIPStr(r))
		}

		// 认证失败，提示 401 Unauthorized
		// Restricted 可以改成其他的值
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		// 401 状态码
		w.WriteHeader(http.StatusUnauthorized)
		util.Log("%q 请求登陆", util.GetRequestIPStr(r))
	}
}
