//go:build !linux

package util

import "net"

func setLinuxBindToDevice(boundDialer *net.Dialer, ifaceName string) {
}
