package config

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/jeessy2/ddns-go/v4/util"
	"gopkg.in/yaml.v3"
)

// Ipv4Reg IPv4正则
const Ipv4Reg = `((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])`

// Ipv6Reg IPv6正则
const Ipv6Reg = `((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))`

// Config 配置
type Config struct {
	Ipv4 struct {
		Enable bool
		// 获取IP类型 url/netInterface
		GetType      string
		URL          string
		NetInterface string
		Domains      []string
	}
	Ipv6 struct {
		Enable bool
		// 获取IP类型 url/netInterface
		GetType      string
		URL          string
		NetInterface string
		IPv6Reg      string // ipv6匹配正则表达式
		Domains      []string
	}
	DNS DNSConfig
	User
	Webhook
	// 禁止公网访问
	NotAllowWanAccess bool
	TTL               string
}

// DNSConfig DNS配置
type DNSConfig struct {
	// 名称。如：alidns,webhook
	Name   string
	ID     string
	Secret string
}

// ConfigCache ConfigCache
type cacheType struct {
	ConfigSingle *Config
	Err          error
	Lock         sync.Mutex
}

var cache = &cacheType{}

// GetConfigCache 获得配置
func GetConfigCache() (conf Config, err error) {
	cache.Lock.Lock()
	defer cache.Lock.Unlock()

	if cache.ConfigSingle != nil {
		return *cache.ConfigSingle, cache.Err
	}

	// init config
	cache.ConfigSingle = &Config{}

	configFilePath := util.GetConfigFilePath()
	_, err = os.Stat(configFilePath)
	if err != nil {
		cache.Err = err
		return *cache.ConfigSingle, err
	}

	byt, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Println("config.yaml读取失败")
		cache.Err = err
		return *cache.ConfigSingle, err
	}

	err = yaml.Unmarshal(byt, cache.ConfigSingle)
	if err != nil {
		log.Println("反序列化配置文件失败", err)
		cache.Err = err
		return *cache.ConfigSingle, err
	}
	// remove err
	cache.Err = nil
	return *cache.ConfigSingle, err
}

// SaveConfig 保存配置
func (conf *Config) SaveConfig() (err error) {
	cache.Lock.Lock()
	defer cache.Lock.Unlock()

	byt, err := yaml.Marshal(conf)
	if err != nil {
		log.Println(err)
		return err
	}

	configFilePath := util.GetConfigFilePath()
	err = ioutil.WriteFile(configFilePath, byt, 0600)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("配置文件已保存在: %s\n", configFilePath)

	// 清空配置缓存
	cache.ConfigSingle = nil

	return
}

// GetIpv4Addr 获得IPv4地址
func (conf *Config) GetIpv4Addr() (result string) {
	// 判断从哪里获取IP
	if conf.Ipv4.GetType == "netInterface" {
		// 从网卡获取IP
		ipv4, _, err := GetNetInterface()
		if err != nil {
			log.Println("从网卡获得IPv4失败!")
			return
		}

		for _, netInterface := range ipv4 {
			if netInterface.Name == conf.Ipv4.NetInterface && len(netInterface.Address) > 0 {
				return netInterface.Address[0]
			}
		}

		log.Println("从网卡中获得IPv4失败! 网卡名: ", conf.Ipv4.NetInterface)
		return
	}

	client := util.CreateNoProxyHTTPClient("tcp4")
	urls := strings.Split(conf.Ipv4.URL, ",")
	for _, url := range urls {
		url = strings.TrimSpace(url)
		resp, err := client.Get(url)
		if err != nil {
			log.Printf("连接失败! <a target='blank' href='%s'>点击查看接口能否返回IPv4地址</a>\n", url)
			log.Printf("错误信息: %s\n", err)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("读取IPv4结果失败! 接口: ", url)
			continue
		}
		comp := regexp.MustCompile(Ipv4Reg)
		result = comp.FindString(string(body))
		if result != "" {
			return
		} else {
			log.Printf("获取IPv4结果失败! 接口: %s ,返回值: %s\n", url, result)
		}
	}

	return
}

// GetIpv6Addr 获得IPv6地址
func (conf *Config) GetIpv6Addr() (result string) {
	// 判断从哪里获取IP
	if conf.Ipv6.GetType == "netInterface" {
		// 从网卡获取IP
		_, ipv6, err := GetNetInterface()
		if err != nil {
			log.Println("从网卡获得IPv6失败!")
			return
		}

		for _, netInterface := range ipv6 {
			if netInterface.Name == conf.Ipv6.NetInterface && len(netInterface.Address) > 0 {
				if conf.Ipv6.IPv6Reg != "" {
					// 匹配第几个IPv6
					if match, err := regexp.MatchString("@\\d", conf.Ipv6.IPv6Reg); err == nil && match {
						num, err := strconv.Atoi(conf.Ipv6.IPv6Reg[1:])
						if err == nil {
							if num > 0 {
								log.Printf("IPv6将使用第 %d 个IPv6地址\n", num)
								if num <= len(netInterface.Address) {
									return netInterface.Address[num-1]
								}
								log.Printf("未找到第 %d 个IPv6地址! 将使用第一个IPv6地址\n", num)
								return netInterface.Address[0]
							}
							log.Printf("IPv6匹配表达式 %s 不正确! 最小从1开始\n", conf.Ipv6.IPv6Reg)
							return
						}
					}
					// 正则表达式匹配
					log.Printf("IPv6将使用正则表达式 %s 进行匹配\n", conf.Ipv6.IPv6Reg)
					for i := 0; i < len(netInterface.Address); i++ {
						matched, err := regexp.MatchString(conf.Ipv6.IPv6Reg, netInterface.Address[i])
						if matched && err == nil {
							log.Println("匹配成功! 匹配到地址: ", netInterface.Address[i])
							return netInterface.Address[i]
						}
						log.Printf("第 %d 个地址 %s 不匹配, 将匹配下一个地址\n", i+1, netInterface.Address[i])
					}
					log.Println("没有匹配到任何一个IPv6地址, 将使用第一个地址")
				}
				return netInterface.Address[0]
			}
		}

		log.Println("从网卡中获得IPv6失败! 网卡名: ", conf.Ipv6.NetInterface)
		return
	}

	client := util.CreateNoProxyHTTPClient("tcp6")
	urls := strings.Split(conf.Ipv6.URL, ",")
	for _, url := range urls {
		url = strings.TrimSpace(url)
		resp, err := client.Get(url)
		if err != nil {
			log.Printf("连接失败! <a target='blank' href='%s'>点击查看接口能否返回IPv6地址</a>, 参考说明:<a target='blank' href='%s'>点击访问</a>\n", url, "https://github.com/jeessy2/ddns-go#使用ipv6")
			log.Printf("错误信息: %s\n", err)
			continue
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("读取IPv6结果失败! 接口: ", url)
			continue
		}
		comp := regexp.MustCompile(Ipv6Reg)
		result = comp.FindString(string(body))
		if result != "" {
			return
		} else {
			log.Printf("获取IPv6结果失败! 接口: %s ,返回值: %s\n", url, result)
		}
	}

	return
}
