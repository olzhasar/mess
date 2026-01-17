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
	gitCmd.Flags().BoolP("verbose", "v", false, "Include diff stats output")
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
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		if !dirty {
			for _, result := range results {
				fmt.Fprintln(cmd.OutOrStdout(), result)
			}
			return
		}

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		repos, err := lib.FindDirtyGitReposWithStats(args[0], older, limitConcurrency, verbose)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		for _, repo := range repos {
			if repo.Err != nil {
				cmd.PrintErrln(repo.Err)
				fmt.Fprintln(cmd.OutOrStdout(), repo.Path)
				continue
			}
			if !verbose {
				fmt.Fprintln(cmd.OutOrStdout(), repo.Path)
				continue
			}
			line := fmt.Sprintf(
				"%s\t%s=%s\t%s/%s",
				repo.Path,
				"files",
				fmt.Sprintf("%d", repo.Stats.ChangedFiles),
				fmt.Sprintf("+%d", repo.Stats.AddedLines),
				fmt.Sprintf("-%d", repo.Stats.DeletedLines),
			)
			if repo.Stats.UntrackedFiles > 0 {
				line += fmt.Sprintf(
					"\t%s=%s",
					"untracked",
					fmt.Sprintf("%d", repo.Stats.UntrackedFiles),
				)
			}
			fmt.Fprintln(cmd.OutOrStdout(), line)
		}
	},
}
