//go:build linux

package util

import (
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

func setLinuxBindToDevice(boundDialer *net.Dialer, ifaceName string) {
	boundDialer.Control = func(network, address string, c syscall.RawConn) error {
		var socketErr error
		err := c.Control(func(fd uintptr) {
			socketErr = unix.SetsockoptString(int(fd), unix.SOL_SOCKET, unix.SO_BINDTODEVICE, ifaceName)
		})
		if err != nil {
			Log("设置 SO_BINDTODEVICE 失败, 回退为仅 LocalAddr 绑定. 网卡: %s, 错误: %v", ifaceName, err)
			return nil
		}
		if socketErr != nil {
			Log("设置 SO_BINDTODEVICE 失败, 回退为仅 LocalAddr 绑定. 网卡: %s, 错误: %v", ifaceName, socketErr)
			return nil
		}
		return nil
	}
}
