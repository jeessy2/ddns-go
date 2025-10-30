//go:build windows

package osutil

import (
	"os"
	"syscall"
)

// StartDetachedProcess starts a process detached from console on Windows.
func StartDetachedProcess(exe string, args []string, nullFile *os.File) (*os.Process, error) {
	const (
		DETACHED_PROCESS         = 0x00000008
		CREATE_NEW_PROCESS_GROUP = 0x00000200
		CREATE_NO_WINDOW         = 0x08000000
	)

	attr := &os.ProcAttr{
		Env:   append(os.Environ(), "DDNS_GO_DAEMON=1"),
		Files: []*os.File{nullFile, nullFile, nullFile},
		Sys:   &syscall.SysProcAttr{CreationFlags: DETACHED_PROCESS | CREATE_NEW_PROCESS_GROUP | CREATE_NO_WINDOW, HideWindow: true},
	}

	return os.StartProcess(exe, args, attr)
}
