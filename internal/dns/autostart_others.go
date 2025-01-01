//go:build !windows

package dns

import "github.com/ducng99/gohole/internal/logger"

func RegisterAutostart() error {
	logger.Printf(logger.LogWarn, "Autostart command is only supported on Windows at the moment.\nYou have to manually set up autostarting gohole on Linux or MacOS.\n")
	return nil
}
