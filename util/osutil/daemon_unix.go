//go:build !windows

package osutil

import (
	"os"
	"syscall"
)

// StartDetachedProcess starts a process detached from terminal on Unix-like systems.
func StartDetachedProcess(exe string, args []string, nullFile *os.File) (*os.Process, error) {
	attr := &os.ProcAttr{
		Env:   append(os.Environ(), "DDNS_GO_DAEMON=1"),
		Files: []*os.File{nullFile, nullFile, nullFile},
		Sys:   &syscall.SysProcAttr{Setsid: true},
	}
	return os.StartProcess(exe, args, attr)
}
