package web

import (
	"ddns-go/config"
	"ddns-go/util"

	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

// Writing 步骤二，填写信息
func Writing(writer http.ResponseWriter, request *http.Request) {
	data, err := Asset("static/pages/writing.html")
	if err != nil {
		// Asset was not found.
	}
	tempFile := os.TempDir() + string(os.PathSeparator) + "writing.html"
	ioutil.WriteFile(tempFile, data, 0600)
	tmpl, err := template.ParseFiles(tempFile)
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
	}
	if conf.Ipv6.URL == "" {
		conf.Ipv6.URL = "https://api-ipv6.ip.sb/ip"
	}

	tmpl.Execute(writer, conf)
}
