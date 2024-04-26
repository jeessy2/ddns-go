package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

//go:embed login.html
var loginEmbedFile embed.FS

// just one token
var token string = ""

// 登录检测
type loginDetect struct {
	failTimes int // 失败次数
}

var ld = &loginDetect{}

// LoginPage login page
func LoginPage(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFS(loginEmbedFile, "login.html")
	if err != nil {
		fmt.Println("Error happened..")
		fmt.Println(err)
		return
	}

	err = tmpl.Execute(writer, struct{}{})
	if err != nil {
		fmt.Println("Error happened..")
		fmt.Println(err)
	}
}

// Login login
func Login(w http.ResponseWriter, r *http.Request) {

	if ld.failTimes > 5 {
		returnError(w, util.LogStr("登录失败次数过多，请稍后再试"))
		return
	}

	// 从请求中读取 JSON 数据
	var data struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		returnError(w, err.Error())
		return
	}

	conf, _ := config.GetConfigCached()

	// 登陆成功
	if data.Username == conf.Username && data.Password == conf.Password {
		ld.failTimes = 0
		token = util.GenerateToken(data.Username)

		// return cookie
		cookie := http.Cookie{
			Name:    "token",
			Value:   token,
			Path:    "/",
			Expires: time.Now().Add(time.Hour * 24), // 设置cookie过期时间为24小时
			Secure:  true,                           // 将Secure设置为true以启用HTTPS安全cookie
		}
		http.SetCookie(w, &cookie)

		util.Log("%q 登陆成功", util.GetRequestIPStr(r))
		returnOK(w, util.LogStr("登陆成功"), token)
		return
	}

	ld.failTimes = ld.failTimes + 1
	util.Log("%q 帐号密码不正确", util.GetRequestIPStr(r))
	returnError(w, util.LogStr("用户名或密码错误"))
}
