//go:build linux

package util

import (
	"net"
	"syscall"
)

func setLinuxBindToDevice(boundDialer *net.Dialer, ifaceName string) {
	// Linux constants for SO_BINDTODEVICE.
	const (
		solSocket      = 1
		soBindToDevice = 25
	)
	boundDialer.Control = func(network, address string, c syscall.RawConn) error {
		var socketErr error
		err := c.Control(func(fd uintptr) {
			socketErr = syscall.SetsockoptString(int(fd), solSocket, soBindToDevice, ifaceName)
		})
		if err != nil {
			Log("设置 SO_BINDTODEVICE 失败, 回退为仅 LocalAddr 绑定. 网卡: %s, 错误: %s", ifaceName, err)
			return nil
		}
		if socketErr != nil {
			Log("设置 SO_BINDTODEVICE 失败, 回退为仅 LocalAddr 绑定. 网卡: %s, 错误: %s", ifaceName, socketErr)
			return nil
		}
		return nil
	}
}
