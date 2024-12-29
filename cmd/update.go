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
package cmd

import (
	"errors"
	"strconv"

	"github.com/ducng99/gohole/cmd/globalFlags"
	"github.com/ducng99/gohole/internal/logger"
	"github.com/ducng99/gohole/internal/sources"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use: "update [sourceID]",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			sourceID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			if sourceID < 1 {
				return errors.New("sourceID must be larger than 0")
			}
		}

		return nil
	},
	Short: "Update all sources or the specified source",
	Run: func(cmd *cobra.Command, args []string) {
		var sourceID int64 = 0

		if len(args) > 0 {
			sourceID, _ = strconv.ParseInt(args[0], 10, 64)
		}

		if sourceID > 0 {
			if err := sources.UpdateSource(sourceID); err != nil {
				if globalFlags.Verbose {
					logger.Printf(logger.LogError, "%v\n", err)
				}
				return
			}
		} else {
			if err := sources.UpdateAllSources(); err != nil {
				if globalFlags.Verbose {
					logger.Printf(logger.LogError, "%v\n", err)
				}
				return
			}
		}

		if err := sources.WriteDomainsToHosts(); err != nil {
			if globalFlags.Verbose {
				logger.Printf(logger.LogError, "%v\n", err)
			}
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
