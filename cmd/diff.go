package cmd

import (
	"dipher/pkg"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "checks are there any breaking changes",
	Long:  "detects breaking changes between two swagger 2.0 specifications",
	Run: func(cmd *cobra.Command, args []string) {
		specV1 := readSpec(args[0])
		specV2 := readSpec(args[1])

		reports := pkg.Diff(specV1, specV2)

		if len(reports) > 0 {
			fmt.Println("breaking changes detected", makeOutput(reports))

			if !safeMode {
				os.Exit(1)
			}
		} else {
			fmt.Println("specs don't have breaking changes")
		}

	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires two swagger specification")
		}

		return nil
	},
}
