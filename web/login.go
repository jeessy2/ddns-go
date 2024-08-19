package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/jeessy2/ddns-go/v6/config"
	"github.com/jeessy2/ddns-go/v6/util"
)

//go:embed login.html
var loginEmbedFile embed.FS

// CookieName cookie name
var cookieName = "token"

// CookieInSystem only one cookie
var cookieInSystem = &http.Cookie{}

// 服务启动时间
var startTime = time.Now()

// 保存限制时间
var saveLimit = time.Duration(30 * time.Minute)

// 登录检测
type loginDetect struct {
	failedTimes uint32       // 失败次数
	ticker      *time.Ticker // 定时器
}

var ld = &loginDetect{ticker: time.NewTicker(5 * time.Minute)}

// Login login page
func Login(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFS(loginEmbedFile, "login.html")
	if err != nil {
		fmt.Println("Error happened..")
		fmt.Println(err)
		return
	}

	conf, _ := config.GetConfigCached()

	err = tmpl.Execute(writer, struct {
		EmptyUser bool // 未填写用户名和密码
	}{
		EmptyUser: conf.Username == "" && conf.Password == "",
	})
	if err != nil {
		fmt.Println("Error happened..")
		fmt.Println(err)
	}
}

// LoginFunc login func
func LoginFunc(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept-Language")
	util.InitLogLang(accept)

	if ld.failedTimes >= 5 {
		lockMinute := loginUnlock()
		returnError(w, util.LogStr("登录失败次数过多，请等待 %d 分钟后再试", lockMinute))
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

	// 用户名密码不能为空
	if data.Username == "" || data.Password == "" {
		returnError(w, util.LogStr("必须输入用户名/密码"))
		return
	}

	conf, _ := config.GetConfigCached()

	// 初始化用户名密码
	if conf.Username == "" && conf.Password == "" {
		if time.Since(startTime) > saveLimit {
			returnError(w, util.LogStr("需在 %s 之前完成用户名密码设置,请重启ddns-go", startTime.Add(saveLimit).Format("2006-01-02 15:04:05")))
			return
		}

		conf.NotAllowWanAccess = true
		u, err := url.Parse(r.Header.Get("referer"))
		if err == nil && !util.IsPrivateNetwork(u.Host) {
			conf.NotAllowWanAccess = false
		}

		conf.Username = data.Username
		hashedPwd, err := conf.CheckPassword(data.Password)
		if err != nil {
			returnError(w, err.Error())
			return
		}
		conf.Password = hashedPwd
		conf.SaveConfig()
	}

	// 登录
	if data.Username == conf.Username && util.PasswordOK(conf.Password, data.Password) {
		ld.ticker.Stop()
		ld.failedTimes = 0

		// 设置cookie过期时间为1天
		timeoutDays := 1
		if conf.NotAllowWanAccess {
			// 内网访问cookie过期时间为30天
			timeoutDays = 30
		}

		// 覆盖cookie
		cookieInSystem = &http.Cookie{
			Name:     cookieName,
			Value:    util.GenerateToken(data.Username), // 生成token
			Path:     "/",
			Expires:  time.Now().AddDate(0, 0, timeoutDays), // 设置过期时间
			HttpOnly: true,
		}
		// 写入cookie
		http.SetCookie(w, cookieInSystem)

		util.Log("%q 登录成功", util.GetRequestIPStr(r))

		returnOK(w, util.LogStr("登录成功"), cookieInSystem.Value)
		return
	}

	ld.failedTimes = ld.failedTimes + 1
	util.Log("%q 帐号密码不正确", util.GetRequestIPStr(r))
	returnError(w, util.LogStr("用户名或密码错误"))
}

// loginUnlock login unlock, return minute
func loginUnlock() (minute uint32) {
	ld.failedTimes = ld.failedTimes + 1
	x := ld.failedTimes
	if x > 1440 {
		x = 1440 // 最多等待一天
	}
	ld.ticker.Reset(time.Duration(x) * time.Minute)

	go func(ticker *time.Ticker) {
		for range ticker.C {
			ld.failedTimes = 4
			ticker.Stop()
		}
	}(ld.ticker)

	return x
}
