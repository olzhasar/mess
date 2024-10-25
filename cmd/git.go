package cmd

import (
	"fmt"

	"github.com/olzhasar/mess/lib"
	"github.com/spf13/cobra"
)

func init() {
	gitCmd.Flags().BoolP("dirty", "d", false, "Filter repositories with uncommitted changes")
	gitCmd.Flags().Int("older", 0, "Filter repositories with last commit older than n days")
	gitCmd.Flags().Int("concurrency", 0, "Limit the number of concurrent checks")
	rootCmd.AddCommand(gitCmd)
}

var gitCmd = &cobra.Command{
	Use:   "git [path]",
	Short: "Find git repositories",
	Long:  `Find git repositories in the specified path`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dirty, err := cmd.Flags().GetBool("dirty")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		older, err := cmd.Flags().GetInt("older")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		limitConcurrency, err := cmd.Flags().GetInt("concurrency")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		results, err := lib.FindGitRepos(args[0], dirty, older, limitConcurrency)

		for _, result := range results {
			fmt.Fprintln(cmd.OutOrStdout(), result)
		}
	},
}
