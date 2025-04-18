package cmd

import (
	"fmt"
	"os"

	"github.com/Koenigseder/badm/internal/filesystem"
	"github.com/Koenigseder/badm/internal/git"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(get)
	get.PersistentFlags().BoolVarP(&overrideExistingFiles, "override", "", false, "override existing files")
}

var get = &cobra.Command{
	Use:   "get",
	Short: "Get a Git repository to use with BADM",
	Long:  `Get a Git repository full of Dotfiles and persist those on the system`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a Git remote URL")
			os.Exit(1)
		}

		repoUrl := args[0]

		git.CloneGitRepository(repoUrl, repoPath)
		filesystem.PersistDotfiles(repoPath, repoName, overrideExistingFiles)
	},
}
