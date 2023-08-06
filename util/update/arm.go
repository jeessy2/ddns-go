// Based on https://github.com/creativeprojects/go-selfupdate/blob/v1.1.1/arm.go

package update

import (
	// unsafe 用于从 runtime 包中获取私有变量
	_ "unsafe"
)

//go:linkname goarm runtime.goarm
var goarm uint8
