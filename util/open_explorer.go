package util

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenExplorer 打开本地浏览器
func OpenExplorer(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	default: // Linux
		// 如果在 Termux 中运行则停止，因为 exec.Command 可能会触发 SIGSYS: bad system call
		// https://github.com/docker/compose/issues/10511#issuecomment-1531453435
		if isTermux() {
			return
		}

		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Please open a browser and visit %s to finish the configuration\n", url)
	} else {
		fmt.Printf("Success to open the browser, please configure in the web page\n")
	}
}
