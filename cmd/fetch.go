package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

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
		persistedFiles := filesystem.PersistDotfiles(repoPath, repoName, overrideExistingFiles)

		// Append persisted files to state
		state.LocalFiles = append(state.LocalFiles, persistedFiles...)

		// Persist the state file
		err := state.WriteStateFile(repoPath)
		if err != nil {
			fmt.Println("Unable writing state file:", err)
			os.Exit(1)
		}

		removeDeadSymlinks()
	},
}

func removeDeadSymlinks() {
	for _, file := range state.LocalFiles {
		_, err := filepath.EvalSymlinks(file)
		if err != nil {
			err = os.Remove(file)
			if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
				fmt.Printf("Unable removing dead symlink %s: %v\n", file, err)
				continue
			}

			// Find index of removed file in state
			index := slices.Index(state.LocalFiles, file)
			if index == -1 {
				continue // File was not found in state
			}

			// Remove file from BADM state
			state.LocalFiles = append(state.LocalFiles[:index], state.LocalFiles[index+1:]...)
		}
	}

	// Persist the state file
	err := state.WriteStateFile(repoPath)
	if err != nil {
		fmt.Println("Unable writing state file:", err)
		os.Exit(1)
	}
}
