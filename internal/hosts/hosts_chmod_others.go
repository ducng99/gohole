//go:build !windows
// +build !windows

package hosts

import (
	"io/fs"
	"os"
)

func chmod(filePath string, permissions fs.FileMode) error {
	return os.Chmod(filePath, permissions)
}
