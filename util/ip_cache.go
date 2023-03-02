package util

const MaxTimes = 5

// IpCache 上次IP缓存
type IpCache struct {
	Addr         string // 缓存地址
	Times        int    // 剩余次数
	FailTimes    int    // 失败次数
	ForceCompare bool   // 是否强制比对
}

var Ipv4Cache *IpCache = &IpCache{}
var Ipv6Cache *IpCache = &IpCache{}

func (d *IpCache) NewIP(newAddr string) {
	if newAddr == "" {
		d.FailTimes++
		d.Times = 0
	} else {
		d.FailTimes = 0
		// 地址改变 或 达到剩余次数 或 强制比对
		if d.Addr != newAddr || d.Times == MaxTimes || d.ForceCompare {
			d.Times = 0
		} else {
			d.Times++
		}
		d.ForceCompare = false
		d.Addr = newAddr
	}
}
