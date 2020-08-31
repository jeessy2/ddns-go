package config

import (
	"ddns-go/util"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

// Ipv4Reg IPV4正则
const Ipv4Reg = `((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])`

// Ipv6Reg IPV6正则
const Ipv6Reg = `((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))`

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
	DNS DNSConfig
}

// DNSConfig DNS配置
type DNSConfig struct {
	Name   string
	ID     string
	Secret string
}

// InitConfigFromFile 获得配置
func (conf *Config) InitConfigFromFile() error {
	configFilePath := util.GetConfigFilePath()
	_, err := os.Stat(configFilePath)
	if err != nil {
		log.Println("没有找到配置文件！请在网页中输入")
		return err
	}
	byt, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Println("config.yaml读取失败")
		return err
	}
	yaml.Unmarshal(byt, conf)
	return nil
}

// GetIpv4Addr 获得IPV4地址
func (conf *Config) GetIpv4Addr() (result string) {
	if conf.Ipv4.Enable {
		resp, err := http.Get(conf.Ipv4.URL)
		if err != nil {
			log.Println("Failed to get ipv4, URL: ", conf.Ipv6.URL)
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("读取IPV4结果失败, URL: ", conf.Ipv4.URL)
			return
		}
		comp := regexp.MustCompile(Ipv4Reg)
		result = comp.FindString(string(body))
	}
	return
}

// GetIpv6Addr 获得IPV6地址
func (conf *Config) GetIpv6Addr() (result string) {
	if conf.Ipv6.Enable {
		resp, err := http.Get(conf.Ipv6.URL)
		if err != nil {
			log.Println("Failed to get ipv6, URL: ", conf.Ipv6.URL)
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("读取IPV6结果失败, URL: ", conf.Ipv6.URL)
			return
		}
		comp := regexp.MustCompile(Ipv6Reg)
		result = comp.FindString(string(body))
	}
	return
}
