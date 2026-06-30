package cmd

import (
	"github.com/olzhasar/mess/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(patternsCmd)
}

var patternsCmd = &cobra.Command{
	Use:   "patterns",
	Short: "Show active cleanup patterns",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		patterns, err := lib.LoadPatterns()
		if err != nil {
			cmd.PrintErrf("Error loading patterns\n%v\n", err)
		}

		for _, pattern := range patterns {
			cmd.Println(pattern)
		}
	},
}
