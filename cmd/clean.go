package cmd

import (
	"github.com/olzhasar/mess/lib"
	"github.com/spf13/cobra"
)

func init() {
	cleanCmd.Flags().BoolP("verbose", "v", false, "Print removed files/directories")
	cleanCmd.Flags().BoolP("recursive", "r", false, "Recursively scan subdirectories")
	cleanCmd.Flags().StringSlice("patterns", []string{}, "Patterns to be removed")
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
		var err error
		options := lib.CleanOptions{}

		options.Verbose, err = cmd.Flags().GetBool("verbose")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		options.Patterns, err = cmd.Flags().GetStringSlice("patterns")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		options.Recursive, err = cmd.Flags().GetBool("recursive")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		n, err := lib.Clean(args[0], options, cmd.OutOrStdout())
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		if n > 0 {
			cmd.Printf("Successfully removed %d items\n", n)
		} else {
			cmd.Println("Nothing found")
		}
	},
}
