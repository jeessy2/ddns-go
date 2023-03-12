package web

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jeessy2/ddns-go/v5/config"
	"github.com/jeessy2/ddns-go/v5/util"
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

		// 配置文件为空, 超过2天禁止从公网访问
		if err != nil && time.Now().Unix()-startTime > 2*24*60*60 &&
			(!util.IsPrivateNetwork(r.RemoteAddr) || !util.IsPrivateNetwork(r.Host)) {
			w.WriteHeader(http.StatusForbidden)
			log.Printf("配置文件为空, 超过2天禁止从公网访问。RemoteAddr: %s\n", r.RemoteAddr)
			return
		}

		// 禁止公网访问
		if conf.NotAllowWanAccess {
			if !util.IsPrivateNetwork(r.RemoteAddr) || !util.IsPrivateNetwork(r.Host) {
				w.WriteHeader(http.StatusForbidden)
				log.Printf("%s 被禁止从公网访问\n", r.RemoteAddr)
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
			log.Printf("%s 登陆失败超过5次! 并延时5分钟响应\n", r.RemoteAddr)
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
			log.Printf("%s 登陆失败!\n", r.RemoteAddr)
		}

		// 认证失败，提示 401 Unauthorized
		// Restricted 可以改成其他的值
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		// 401 状态码
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("%s 请求登陆!\n", r.RemoteAddr)
	}
}
