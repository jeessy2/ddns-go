package config

import (
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/jeessy2/ddns-go/v4/util"
	"gopkg.in/yaml.v3"
)

// Ipv4Reg IPv4正则
var Ipv4Reg = regexp.MustCompile(`((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])`)

// Ipv6Reg IPv6正则
var Ipv6Reg = regexp.MustCompile(`((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))`)

// Config 配置
type Config struct {
	Ipv4 struct {
		Enable bool
		// 获取IP类型 url/netInterface
		GetType      string
		URL          string
		NetInterface string
		Cmd          string
		Domains      []string
	}
	Ipv6 struct {
		Enable bool
		// 获取IP类型 url/netInterface
		GetType      string
		URL          string
		NetInterface string
		Cmd          string
		IPv6Reg      string // ipv6匹配正则表达式
		Domains      []string
	}
	DNS DNSConfig
	TTL string
}

// DNSConfig DNS配置
type DNSConfig struct {
	// 名称。如：alidns,webhook
	ID     string
	Secret string
}

type ConfigGlobal struct {
	User
	Webhook
	// 禁止公网访问
	NotAllowWanAccess bool
}

type ConfigFile struct {
	ConfigMap    map[string]Config
	ConfigGlobal ConfigGlobal
}

// ConfigCache ConfigCache
type cacheType struct {
	configMap    *map[string]Config
	configGlobal *ConfigGlobal
	Err          error
	Lock         sync.Mutex
}

var cache = &cacheType{}

// GetConfigCache 获得配置
func GetConfigGlobal() (conf ConfigGlobal, err error) {
	cache.Lock.Lock()
	defer cache.Lock.Unlock()

	if cache.configGlobal != nil {
		return *cache.configGlobal, cache.Err
	}

	// init config
	cache.configGlobal = &ConfigGlobal{}

	configFilePath := util.GetConfigFilePath()
	_, err = os.Stat(configFilePath)
	if err != nil {
		cache.Err = err
		return *cache.configGlobal, err
	}

	byt, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Println("config.yaml读取失败")
		cache.Err = err
		return *cache.configGlobal, err
	}

	configFile := &ConfigFile{}
	err = yaml.Unmarshal(byt, configFile)
	if err != nil {
		log.Println("反序列化配置文件失败", err)
		cache.Err = err
		return *cache.configGlobal, err
	}
	if configFile.ConfigGlobal.Username == "" && configFile.ConfigGlobal.Password == "" {
		configFile.ConfigGlobal.NotAllowWanAccess = true
	}
	cache.configGlobal = &configFile.ConfigGlobal
	if len(configFile.ConfigMap) == 0 {
		cache.configMap = &map[string]Config{}
	} else {
		cache.configMap = &configFile.ConfigMap
	}
	// remove err
	cache.Err = nil
	return *cache.configGlobal, err
}

func GetConfigMap() (conf map[string]Config) {
	if cache.configMap == nil {
		cache.configMap = &map[string]Config{}
		GetConfigGlobal()
	}
	return *cache.configMap
}

// SaveConfig 保存配置
func SaveConfig(cglobal ConfigGlobal, cmap map[string]Config) (err error) {
	cache.Lock.Lock()
	defer cache.Lock.Unlock()

	configFile := ConfigFile{ConfigMap: cmap, ConfigGlobal: cglobal}
	byt, err := yaml.Marshal(configFile)
	if err != nil {
		log.Println(err)
		return err
	}

	configFilePath := util.GetConfigFilePath()
	err = os.WriteFile(configFilePath, byt, 0600)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("配置文件已保存在: %s\n", configFilePath)

	// 清空配置缓存
	cache.configGlobal = nil
	cache.configMap = nil

	return
}

func CompatibleConfig() {
	conf, err := GetConfigGlobal()
	if err != nil || conf.NotAllowWanAccess || conf.Password != "" {
		return
	}
	temp := map[interface{}]interface{}{}
	byt, err := os.ReadFile(util.GetConfigFilePath())
	if err == nil {
		err = yaml.Unmarshal(byt, temp)
	}
	if err != nil {
		return
	}

	wan := temp["notallowwanaccess"]
	switch wan := wan.(type) {
	case bool:
		conf.NotAllowWanAccess = wan
	}
	user := temp["user"]
	switch user := user.(type) {
	case map[string]interface{}:
		username := user["username"]
		switch username := username.(type) {
		case string:
			conf.Username = username
		}
		password := user["password"]
		switch password := password.(type) {
		case string:
			conf.Password = password
		}
	}
	hook := temp["webhook"]
	switch hook := hook.(type) {
	case map[string]interface{}:
		url := hook["webhookurl"]
		switch url := url.(type) {
		case string:
			conf.WebhookURL = url
		}
		body := hook["webhookrequestbody"]
		switch body := body.(type) {
		case string:
			conf.WebhookRequestBody = body
		}
	}

	cmap := GetConfigMap()
	dns := temp["dns"]
	switch dns := dns.(type) {
	case map[string]interface{}:
		name := dns["name"]
		switch name := name.(type) {
		case string:
			cs := Config{}
			id := dns["id"]
			switch id := id.(type) {
			case string:
				cs.DNS.ID = id
			}
			secret := dns["secret"]
			switch secret := secret.(type) {
			case string:
				cs.DNS.Secret = secret
			}

			ipv4 := temp["ipv4"]
			switch ipv4 := ipv4.(type) {
			case map[string]interface{}:
				enable := ipv4["enable"]
				switch enable := enable.(type) {
				case bool:
					cs.Ipv4.Enable = enable
				}
				gettype := ipv4["gettype"]
				switch gettype := gettype.(type) {
				case string:
					cs.Ipv4.GetType = gettype
				}
				url := ipv4["url"]
				switch url := url.(type) {
				case string:
					cs.Ipv4.URL = url
				}
				net := ipv4["netinterface"]
				switch net := net.(type) {
				case string:
					cs.Ipv4.NetInterface = net
				}
				cmd := ipv4["cmd"]
				switch cmd := cmd.(type) {
				case string:
					cs.Ipv4.Cmd = cmd
				}
				domains := ipv4["domains"]
				switch domains := domains.(type) {
				case []interface{}:
					sl := []string{}
					for _, v := range domains {
						switch v := v.(type) {
						case string:
							sl = append(sl, v)
						}
					}
					cs.Ipv4.Domains = sl
				}
			}

			ipv6 := temp["ipv6"]
			switch ipv6 := ipv6.(type) {
			case map[string]interface{}:
				enable := ipv6["enable"]
				switch enable := enable.(type) {
				case bool:
					cs.Ipv6.Enable = enable
				}
				gettype := ipv6["gettype"]
				switch gettype := gettype.(type) {
				case string:
					cs.Ipv6.GetType = gettype
				}
				url := ipv6["url"]
				switch url := url.(type) {
				case string:
					cs.Ipv6.URL = url
				}
				net := ipv6["netinterface"]
				switch net := net.(type) {
				case string:
					cs.Ipv6.NetInterface = net
				}
				cmd := ipv6["cmd"]
				switch cmd := cmd.(type) {
				case string:
					cs.Ipv6.Cmd = cmd
				}
				ipreg := ipv6["ipv6reg"]
				switch ipreg := ipreg.(type) {
				case string:
					cs.Ipv6.IPv6Reg = ipreg
				}
				domains := ipv6["domains"]
				switch domains := domains.(type) {
				case []interface{}:
					sl := []string{}
					for _, v := range domains {
						switch v := v.(type) {
						case string:
							sl = append(sl, v)
						}
					}
					cs.Ipv6.Domains = sl
				}
			}

			ttl := temp["ttl"]
			switch ttl := ttl.(type) {
			case string:
				cs.TTL = ttl
			}
			cmap[name] = cs
		}
	}

	cache.configGlobal = &conf
	cache.configMap = &cmap
}

func (conf *Config) getIpv4AddrFromInterface() string {
	ipv4, _, err := GetNetInterface()
	if err != nil {
		log.Println("从网卡获得IPv4失败!")
		return ""
	}

	for _, netInterface := range ipv4 {
		if netInterface.Name == conf.Ipv4.NetInterface && len(netInterface.Address) > 0 {
			return netInterface.Address[0]
		}
	}

	log.Println("从网卡中获得IPv4失败! 网卡名: ", conf.Ipv4.NetInterface)
	return ""
}

func (conf *Config) getIpv4AddrFromUrl() string {
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
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("读取IPv4结果失败! 接口: ", url)
			continue
		}
		result := Ipv4Reg.FindString(string(body))
		if result == "" {
			log.Printf("获取IPv4结果失败! 接口: %s ,返回值: %s\n", url, result)
		}
		return result
	}
	return ""
}

func (conf *Config) getAddrFromCmd(addrType string) string {
	var cmd string
	var comp *regexp.Regexp
	if addrType == "IPv4" {
		cmd = conf.Ipv4.Cmd
		comp = Ipv4Reg
	} else {
		cmd = conf.Ipv6.Cmd
		comp = Ipv6Reg
	}
	// cmd is empty
	if cmd == "" {
		return ""
	}
	// run cmd with proper shell
	var execCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		execCmd = exec.Command("powershell", "-Command", cmd)
	} else {
		execCmd = exec.Command("bash", "-rc", cmd)
	}
	// run cmd
	out, err := execCmd.Output()
	if err != nil {
		log.Printf("获取%s结果失败! 未能成功执行命令：%s，错误：%q\n", addrType, execCmd.String(), err.Error())
		return ""
	}
	str := string(out)
	// get result
	result := comp.FindString(str)
	if result == "" {
		log.Printf("获取%s结果失败! 命令：%s，标准输出：%q\n", addrType, execCmd.String(), str)
	}
	return result
}

// GetIpv4Addr 获得IPv4地址
func (conf *Config) GetIpv4Addr() string {
	// 判断从哪里获取IP
	switch conf.Ipv4.GetType {
	case "netInterface":
		// 从网卡获取 IP
		return conf.getIpv4AddrFromInterface()
	case "url":
		// 从 URL 获取 IP
		return conf.getIpv4AddrFromUrl()
	case "cmd":
		// 从命令行获取 IP
		return conf.getAddrFromCmd("IPv4")
	}
	return "" // unknown type
}

func (conf *Config) getIpv6AddrFromInterface() string {
	_, ipv6, err := GetNetInterface()
	if err != nil {
		log.Println("从网卡获得IPv6失败!")
		return ""
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
						return ""
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
	return ""
}

func (conf *Config) getIpv6AddrFromUrl() string {
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
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("读取IPv6结果失败! 接口: ", url)
			continue
		}
		result := Ipv6Reg.FindString(string(body))
		if result == "" {
			log.Printf("获取IPv6结果失败! 接口: %s ,返回值: %s\n", url, result)
		}
		return result
	}
	return ""
}

// GetIpv6Addr 获得IPv6地址
func (conf *Config) GetIpv6Addr() (result string) {
	// 判断从哪里获取IP
	switch conf.Ipv6.GetType {
	case "netInterface":
		// 从网卡获取 IP
		return conf.getIpv6AddrFromInterface()
	case "url":
		// 从 URL 获取 IP
		return conf.getIpv6AddrFromUrl()
	case "cmd":
		// 从命令行获取 IP
		return conf.getAddrFromCmd("IPv6")
	}
	return "" // unknown type
}
