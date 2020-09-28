package util

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenExplorer 打开本地浏览器
func OpenExplorer(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		// mac
		cmd = "open"
	default:
		// linux
		cmd = "xdg-open"
	}
	args = append(args, url)

	err := exec.Command(cmd, args...).Start()
	if err != nil {
		fmt.Printf("自动打开浏览器失败, 请手动打开。错误信息: %s", err.Error())
	} else {
		fmt.Println("成功打开浏览器")
	}
}
