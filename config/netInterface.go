package config

import (
	"fmt"
	"net"
)

// NetInterface 本机网络
type NetInterface struct {
	Name    string
	Address []string
}

// GetNetInterface 获得网卡地址
// 返回ipv4, ipv6地址
func GetNetInterface() (ipv4NetInterfaces []NetInterface, ipv6NetInterfaces []NetInterface, err error) {
	allNetInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
		return ipv4NetInterfaces, ipv6NetInterfaces, err
	}

	// https://en.wikipedia.org/wiki/IPv6_address#General_allocation
	_, ipv6Unicast, _ := net.ParseCIDR("2000::/3")

	for i := 0; i < len(allNetInterfaces); i++ {
		if (allNetInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := allNetInterfaces[i].Addrs()
			ipv4 := []string{}
			ipv6 := []string{}

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
					_, bits := ipnet.Mask.Size()
					// 需匹配全局单播地址
					if bits == 128 && ipv6Unicast.Contains(ipnet.IP) {
						ipv6 = append(ipv6, ipnet.IP.String())
					}
					if bits == 32 {
						ipv4 = append(ipv4, ipnet.IP.String())
					}
				}
			}

			if len(ipv4) > 0 {
				ipv4NetInterfaces = append(
					ipv4NetInterfaces,
					NetInterface{
						Name:    allNetInterfaces[i].Name,
						Address: ipv4,
					},
				)
			}

			if len(ipv6) > 0 {
				ipv6NetInterfaces = append(
					ipv6NetInterfaces,
					NetInterface{
						Name:    allNetInterfaces[i].Name,
						Address: ipv6,
					},
				)
			}

		}
	}

	return ipv4NetInterfaces, ipv6NetInterfaces, nil
}
