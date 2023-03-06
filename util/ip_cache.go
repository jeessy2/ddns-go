package util

// IpCache 上次IP缓存
type IpCache struct {
	Addr      string // 缓存地址
	Times     int    // 剩余次数
	FailTimes int    // 获取ip失败的次数
}

var ForceCompare = true

func (d *IpCache) Check(newAddr string) bool {
	if newAddr == "" {
		return true
	}
	// 地址改变 或 达到剩余次数 或 强制比对
	if d.Addr != newAddr || d.Times <= 1 {
		d.Addr = newAddr
		d.Times = 6
		return true
	}
	d.Addr = newAddr
	d.Times--
	return false
}
