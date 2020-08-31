package web

import (
	"ddns-go/config"
	"ddns-go/util"
	"log"

	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

// Writing 步骤二，填写信息
func Writing(writer http.ResponseWriter, request *http.Request) {
	tempPath, err := util.GetStaticResourcePath("static/pages/writing.html")
	if err != nil {
		log.Println("Asset was not found.")
		return
	}
	tmpl, err := template.ParseFiles(tempPath)
	if err != nil {
		fmt.Println("Error happened..")
		fmt.Println(err)
		return
	}

	conf := &config.Config{}

	// 解析文件
	var configFile string = util.GetConfigFilePath()
	_, err = os.Stat(configFile)
	if err == nil {
		// 不为空，解析文件
		byt, err := ioutil.ReadFile(configFile)
		if err == nil {
			err = yaml.Unmarshal(byt, conf)
			if err == nil {
				tmpl.Execute(writer, conf)
				return
			}
		}
	}

	// 默认值
	if conf.Ipv4.URL == "" {
		conf.Ipv4.URL = "https://api-ipv4.ip.sb/ip"
		conf.Ipv4.Enable = true
	}
	if conf.Ipv6.URL == "" {
		conf.Ipv6.URL = "https://api-ipv6.ip.sb/ip"
	}
	if conf.DNS.Name == "" {
		conf.DNS.Name = "alidns"
	}

	tmpl.Execute(writer, conf)
}
