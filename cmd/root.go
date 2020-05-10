package cmd

import (
	"differ/internal/helper"
	"differ/pkg/differ"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "differ",
	Short: "Swagger diff for breaking changes detect",
	Run: func(cmd *cobra.Command, args []string) {
		specV1 := helper.ReadSpec(args[0])
		specV2 := helper.ReadSpec(args[1])

		diff := differ.Diff(specV1, specV2)

		if len(diff) > 0 {
			log.Fatalln("breaking changes detected", helper.MakeReport(diff))
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
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of differ",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0.0.1")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
