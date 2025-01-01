//go:build windows

package hosts

import (
	"path/filepath"

	"github.com/ducng99/gohole/internal/logger"
	"golang.org/x/sys/windows"
)

func getHostsFilePath() (string, error) {
	windowsDir, err := windows.GetWindowsDirectory()
	if err != nil {
		logger.Printf(logger.LogError, "Cannot get Windows directory\n")
		return "", err
	}

	return filepath.Join(windowsDir, "System32", "drivers", "etc", "hosts"), nil
}
