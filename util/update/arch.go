// Based on https://github.com/creativeprojects/go-selfupdate/blob/v1.1.1/arch.go

package update

import (
	"fmt"
	"runtime"
)

const (
	minARM = 5
	maxARM = 7
)

// generateAdditionalArch 可以根据 CPU 类型使用
func generateAdditionalArch() []string {
	if runtime.GOARCH == "arm" && goarm >= minARM && goarm <= maxARM {
		additionalArch := make([]string, 0, maxARM-minARM)
		for v := goarm; v >= minARM; v-- {
			additionalArch = append(additionalArch, fmt.Sprintf("armv%d", v))
		}
		return additionalArch
	}
	if runtime.GOARCH == "amd64" {
		return []string{"x86_64"}
	}
	return []string{}
}
