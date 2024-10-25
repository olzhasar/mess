package cmd

import (
	"github.com/olzhasar/mess/lib"
	"github.com/spf13/cobra"
)

func init() {
	cleanCmd.Flags().BoolP("verbose", "v", false, "Print removed files/directories")
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean [path]",
	Short: "Delete temporary files",
	Long: `Delete temporary files in the specified path

Patterns:
python: *.pyc, __pycache__, .mypy_cache, .pytest_cache, .ruff_cache, .tox, .nox
node: node_modules
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		n, err := lib.Clean(args[0], verbose)
		if err != nil {
			cmd.PrintErrln(err)
		}

		cmd.Printf("Successfully removed %d items\n", n)
	},
}
