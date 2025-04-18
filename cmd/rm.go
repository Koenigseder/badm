package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Koenigseder/badm/internal/filesystem"
	"github.com/Koenigseder/badm/internal/git"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rm)
}

var rm = &cobra.Command{
	Use:   "rm",
	Short: "Remove a file (or files) from your Dotfiles",
	Long:  `Remove a file (or files) from your Dotfiles. The file does not get deleted, it gets placed at the original location`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Fetch Git remote status and persist all Dotfiles to the system using symlinks
		fmt.Println("Fetching remote status...")

		if git.FetchAndUpdate(repoPath) {
			fmt.Println("Persisting changes...")
			filesystem.PersistDotfiles(repoPath, repoName, overrideExistingFiles)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		removeDotfiles(args)
	},
}

// args[] contains all relative (or absolute) paths to Dotfiles which should be restored
func removeDotfiles(args []string) {
	if len(args) == 0 {
		fmt.Println("No Dotfiles specified")
		os.Exit(1)
	}

	// Get current working directory
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable retrieving current directory path:", err)
		os.Exit(1)
	}

	// Remove all Dotfiles
	for _, relativeFilePath := range args {
		var longFilePath string

		// Use the passed file path if it begins with '~' or '/'
		if strings.HasPrefix(relativeFilePath, "~") || strings.HasPrefix(relativeFilePath, "/") {
			longFilePath = relativeFilePath
		} else {
			// Not optimized path to file
			longFilePath = fmt.Sprintf("%s/%s", pwd, relativeFilePath)
		}

		// Clean path
		shortAbsoluteFilePath := path.Clean(longFilePath)

		// Intercept BADM repo path to get the file from there
		fileRepoPath := path.Join(repoPath, strings.Replace(shortAbsoluteFilePath, homeDir, "", 1))

		fmt.Printf("Removing %s from Dotfiles...\n", relativeFilePath)

		// Remove soft symbolic link
		err = os.Remove(shortAbsoluteFilePath)
		if err != nil {
			fmt.Printf("Unable removing soft symbolic link at %s: %s\n", shortAbsoluteFilePath, err)
			os.Exit(1)
		}

		// Move file
		err = os.Rename(fileRepoPath, shortAbsoluteFilePath)
		if err != nil {
			fmt.Println("Unable moving file:", err)
			os.Exit(1)
		}
	}

	// Remove Dotfiles from Git and push it
	git.CommitAndPushFiles(repoPath, "💥 Remove file")

	fmt.Printf("Successfully removed and pushed Dotfiles!")
}
