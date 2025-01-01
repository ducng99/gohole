//go:build windows

package upgrader

import (
	"os/exec"
)

func pkill(process string) error {
	killCmd := exec.Command("taskkill", "/IM", process, "/T", "/F")
	killCmd.Stdout = nil
	killCmd.Stderr = nil

	return killCmd.Run()
}
