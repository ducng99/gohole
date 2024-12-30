//go:build !windows
// +build !windows

package dns

func RegisterAutostart() error {
	return nil
}
