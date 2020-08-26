package config

import (
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

// Config 配置
type Config struct {
	Ipv4 struct {
		Enable  bool
		URL     string
		Domains []string
	}
	Ipv6 struct {
		Enable  bool
		URL     string
		Domains []string
	}
	DNS struct {
		Name   string
		ID     string
		Secret string
	}
}

// GetConfigFromFile 获得配置
func (conf *Config) GetConfigFromFile() {
	byt, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Println("config.yaml读取失败")
	}
	yaml.Unmarshal(byt, conf)
}

// GetIpv4Addr 获得IPV4地址
func (conf *Config) GetIpv4Addr() (result string, err error) {
	resp, err := http.Get(conf.Ipv4.URL)
	if err != nil {
		log.Println("获得IPV4失败")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取IPV4结果失败")
		return
	}
	result = string(body)
	result = "8.8.8.8"
	return
}

// GetIpv6Addr 获得IPV6地址
func (conf *Config) GetIpv6Addr() (result string, err error) {
	resp, err := http.Get(conf.Ipv6.URL)
	if err != nil {
		log.Println("获得IPV6失败")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取IPV6结果失败")
		return
	}
	result = string(body)
	return
}
