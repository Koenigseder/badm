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
	rootCmd.AddCommand(add)
	add.PersistentFlags().BoolVarP(&overrideExistingFiles, "override", "", false, "override existing files")
}

var add = &cobra.Command{
	Use:   "add",
	Short: "Add a file (or files) to your Dotfiles",
	Long:  `Add a file (or files) to your Dotfiles. Those files get automatically pushed to your Git repository`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Fetch Git remote status and persist all Dotfiles to the system using symlinks
		fmt.Println("Fetching remote status...")

		if git.FetchAndUpdate(repoPath) {
			fmt.Println("Persisting changes...")
			filesystem.PersistDotfiles(repoPath, repoName, overrideExistingFiles)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		addDotfiles(args)
	},
}

// args[] contains all relative (or absolute) paths to Dotfiles which should be added
func addDotfiles(args []string) {
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

	// Add all Dotfiles
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

		// Intercept BADM repo path to save the file there
		fileRepoPath := path.Join(repoPath, strings.Replace(shortAbsoluteFilePath, homeDir, "", 1))

		if info, _ := os.Lstat(shortAbsoluteFilePath); info.Mode().Type() == os.ModeSymlink {
			fmt.Println(shortAbsoluteFilePath, "is a symlink and cannot be added to the Dotfiles")
			os.Exit(1)
		}

		fmt.Printf("Adding %s to Dotfiles...\n", relativeFilePath)

		// Move the Dotfile to the BADM repo and create symbolic link
		filesystem.MoveFileWithSymLink(shortAbsoluteFilePath, fileRepoPath)
	}

	// Add Dotfile to Git and push it
	git.CommitAndPushFiles(repoPath, "🚀 Add new file")

	fmt.Println("Successfully added and pushed Dotfiles!")
}
