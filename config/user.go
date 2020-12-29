package config

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

// User 登录用户
type User struct {
	Username string
	Password string
}

// ViewFunc func
type ViewFunc func(http.ResponseWriter, *http.Request)

// BasicAuth basic auth
func BasicAuth(f ViewFunc) ViewFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 帐号或密码为空。跳过
		conf, _ := GetConfigCache()
		if conf.Username == "" && conf.Password == "" {
			// 执行被装饰的函数
			f(w, r)
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
					// 执行被装饰的函数
					f(w, r)
					return
				}
			}
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
