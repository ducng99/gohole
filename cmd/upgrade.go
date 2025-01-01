/*
Copyright © 2025 Thomas Nguyen <tom@tomng.dev>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/ducng99/gohole/globals"
	"github.com/ducng99/gohole/internal/logger"
	"github.com/ducng99/gohole/internal/upgrader"
	"github.com/spf13/cobra"
)

var originalFilePath string

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade gohole to latest version",
	Long:  `Check and self upgrade gohole binary to latest binary from GitHub.`,
	Run: func(cmd *cobra.Command, args []string) {
		if originalFilePath == "" {
			if err := upgrader.RunTemp(); err != nil && globals.Verbose {
				logger.Fatalf("An error occurred while self upgrading: %v\n", err)
			}
		} else {
			if err := upgrader.CheckAndUpgrade(originalFilePath); err != nil && globals.Verbose {
				logger.Fatalf("An error occurred while self upgrading: %v\n", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	upgradeCmd.Flags().StringVar(&originalFilePath, "file-path", "", "")
	upgradeCmd.Flags().MarkHidden("file-path")
}
