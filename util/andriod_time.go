package util

import "time"

func FixTimezone(timeZoneStr string) {
	offset, err := getTimeZoneOffset(timeZoneStr)

	// 如果解析不出来, 不修改时区
	if err != nil {
		return
	}

	timeZone := time.FixedZone("GMT", offset*3600)
	time.Local = timeZone
}

// 将时区字符串转换为偏移量（以小时为单位）
func getTimeZoneOffset(timeZoneStr string) (int, error) {
	// 获取时区信息
	location, err := time.LoadLocation(timeZoneStr)
	if err != nil {
		return 0, err
	}

	// 获取当前时间
	currentTime := time.Now().In(location)

	_, offsetInSeconds := currentTime.Zone()

	// 将秒转换为小时
	offsetInHours := offsetInSeconds / 3600

	return offsetInHours, nil
}
