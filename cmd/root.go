package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mess",
	Short: "Clean temporary development files",
	Long:  "mess removes common temporary development files, caches, and build artifacts from a directory.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
