//go:build !windows

package hosts

func getHostsFilePath() (string, error) {
	return "/etc/hosts", nil
}
