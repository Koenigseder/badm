package cmd

import (
	"github.com/Koenigseder/badm/internal/filesystem"
	"github.com/Koenigseder/badm/internal/git"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fetch)
	fetch.PersistentFlags().BoolVarP(&overrideExistingFiles, "override", "", false, "override existing files")
}

var fetch = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch a Git repository's status",
	Long:  `Fetch a Git repository's status and persist Dotfiles on the system`,
	Run: func(cmd *cobra.Command, args []string) {
		git.FetchAndUpdate(repoPath)
		filesystem.PersistDotfiles(repoPath, repoName, overrideExistingFiles)
	},
}
