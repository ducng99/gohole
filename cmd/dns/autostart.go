/*
Copyright © 2024 Thomas Nguyen <tom@tomng.dev>

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
package dns

import (
	"github.com/ducng99/gohole/globals"
	"github.com/ducng99/gohole/internal/dns"
	"github.com/ducng99/gohole/internal/logger"
	"github.com/spf13/cobra"
)

// autoStartCmd represents the start command
var autoStartCmd = &cobra.Command{
	Use:   "autostart",
	Short: "Register gohole to autostart",
	Run: func(cmd *cobra.Command, args []string) {
		if err := dns.RegisterAutostart(); err != nil {
			if globals.Verbose {
				logger.Printf(logger.LogError, "Failed when registering DNS server for autostart: %v\n", err)
			}
			return
		}
	},
}

func init() {
	DnsCmd.AddCommand(autoStartCmd)
}
