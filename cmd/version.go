package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of dipher",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.1.1")
	},
}
