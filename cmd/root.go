package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var safeMode bool

var rootCmd = &cobra.Command{
	Use:   "dipher",
	Short: "Swagger diff for breaking changes detection",
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().BoolVarP(&safeMode,
		"safe-mode", "s",
		false, "In such mode command won't fail, just print breaking changes")
}

// Execute is CLI entry poit
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
