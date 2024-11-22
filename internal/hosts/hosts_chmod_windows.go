//go:build windows
// +build windows

package hosts

import (
	"io/fs"

	acl "github.com/hectane/go-acl"
)

func chmod(filePath string, permissions fs.FileMode) error {
	return acl.Chmod(filePath, permissions)
}
