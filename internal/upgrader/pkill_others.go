//go:build !windows

package upgrader

import "os/exec"

func pkill(process string) error {
	killCmd := exec.Command("pkill", "-9", process)
	killCmd.Stdout = nil
	killCmd.Stderr = nil

	return killCmd.Run()
}
