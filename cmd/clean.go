package cmd

import (
	"fmt"

	"github.com/olzhasar/mess/lib"
	"github.com/spf13/cobra"
)

func init() {
	cleanCmd.Flags().BoolP("verbose", "v", false, "Print each removed path")
	cleanCmd.Flags().BoolP("recursive", "r", false, "Scan subdirectories recursively")
	cleanCmd.Flags().Bool("no-size", false, "Skip freed-space calculation (faster for large directories)")
	cleanCmd.Flags().StringSlice("patterns", []string{}, "Comma-separated patterns to remove instead of configured patterns")
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean [path]",
	Short: "Delete temporary development files",
	Long: `Delete temporary development files from path.

By default, mess checks only direct children of path. Use --recursive or -r to scan subdirectories.

	Built-in patterns are used unless a patterns file exists or --patterns is provided. To check which patterns are currently configured, run: mess patterns`,
	Example: `  mess clean ~/dev/project
  mess clean -r ~/dev
  mess clean -r --patterns "node_modules,*.pyc" ~/dev`,
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

		noSize, err := cmd.Flags().GetBool("no-size")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		options.CalcFreed = !noSize

		result, err := lib.Clean(args[0], options, cmd.OutOrStdout())
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		if result.Count > 0 {
			cmd.Printf("Successfully removed %d item(s)\n", result.Count)
			if options.CalcFreed {
				cmd.Printf("%s freed\n", formatBytes(result.BytesFreed))
			}
		} else {
			cmd.Println("Nothing found")
		}
	},
}

func formatBytes(bytes uint64) string {
	const step = 1024

	if bytes < step {
		return fmt.Sprintf("%d B", bytes)
	}

	result := float64(bytes)

	for _, unit := range []string{"K", "M", "G", "T", "P"} {
		result /= float64(step)
		if result < step {
			return fmt.Sprintf("%.1f%s", result, unit)
		}
	}

	return fmt.Sprintf("%.1f%s", result, "PB")
}
