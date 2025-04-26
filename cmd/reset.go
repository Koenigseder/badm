package cmd

import (
	"github.com/Koenigseder/badm/internal/filesystem"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(reset)
	reset.PersistentFlags().BoolVarP(&dryRun, "dryRun", "", false, "do a dry run on what would happen")
}

var reset = &cobra.Command{
	Use:   "reset",
	Short: "Repopulate all original Dotfiles on the system",
	Long:  `Undo any changes made to the system by BADM. Removes symlinks and repopulates all original Dotfiles on the system`,
	Run: func(cmd *cobra.Command, args []string) {
		filesystem.RestoreDotfiles(repoPath, repoName, dryRun)
	},
}
