package util

import "os"

// isTermux 是否在 Termux 中运行
//
// https://wiki.termux.com/wiki/Getting_started
func isTermux() bool {
	return os.Getenv("PREFIX") == "/data/data/com.termux/files/usr"
}
