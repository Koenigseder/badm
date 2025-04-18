package cmd

import (
	"fmt"

	"github.com/Koenigseder/badm/internal/git"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(save)
}

var save = &cobra.Command{
	Use:   "save",
	Short: "Save the current Dotfiles state",
	Long:  `Save the current Dotfiles state to the remote Git repository`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Saving current Dotfiles state to remote Git repository...")

		git.CommitAndPushFiles(repoPath, "💾 Save files")

		fmt.Println("Saved current Dotfiles state to remote Git repository")
	},
}
