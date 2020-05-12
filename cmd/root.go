package cmd

import (
	"differ/pkg/differ"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var safeMode bool

var rootCmd = &cobra.Command{
	Use:   "differ",
	Short: "Swagger diff for breaking changes detect",
	Run: func(cmd *cobra.Command, args []string) {
		specV1 := readSpec(args[0])
		specV2 := readSpec(args[1])

		diff := differ.Diff(specV1, specV2)

		if len(diff) > 0 {
			if safeMode {
				log.Println("breaking changes detected", makeReport(diff))
			} else {
				log.Fatalln("breaking changes detected", makeReport(diff))
			}
		} else {
			log.Println("specs don't have breaking changes")
		}

	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires two swagger specification")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().BoolVarP(&safeMode, "safe-mode",
		"s", false, "In such mode command won't fail, just print breaking changes")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of differ",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.0.1")
	},
}

// Execute is CLI entry poit
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
