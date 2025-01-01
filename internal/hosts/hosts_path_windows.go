//go:build windows
// +build windows

package hosts

import (
	"path/filepath"
	"syscall"
	"unsafe"

	"github.com/ducng99/gohole/internal/logger"
)

func getHostsFilePath() (string, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getWindowsDirectoryHandle := kernel32.NewProc("GetWindowsDirectoryA")

	windowsDir := make([]byte, 260)

	bytesWritten, _, err := getWindowsDirectoryHandle.Call(uintptr(unsafe.Pointer(&windowsDir[0])), 260)
	if bytesWritten == 0 {
		logger.Printf(logger.LogError, "Cannot get Windows directory\n")
		return "", err
	}

	if uint32(bytesWritten) == 260 {
		logger.Printf(logger.LogWarn, "The Windows directory path has reached maximum length. This could be an error.\n")
	}

	return filepath.Join(string(windowsDir[:bytesWritten]), "System32", "drivers", "etc", "hosts"), nil
}
