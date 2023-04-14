// based on https://gist.github.com/hyg/9c4afcd91fe24316cbf0

package util

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenExplorer 打开本地浏览器
func OpenExplorer(url string) {
	var err error

	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin": // macOS
		err = exec.Command("open", url).Start()
	default: // Linux
		err = exec.Command("xdg-open", url).Start()
	}

	if err != nil {
		fmt.Printf("请手动打开浏览器并访问 %s 进行配置\n", url)
	} else {
		fmt.Println("成功打开浏览器, 请在网页中进行配置")
	}
}
