package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/Koenigseder/badm/internal/git"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new BADM repo",
	Long:  `Create a new BADM repo to manage your Dotfiles`,
	Run: func(_ *cobra.Command, args []string) {
		createNewBadmRepo(args)
	},
}

func createNewBadmRepo(args []string) {
	if len(args) == 0 {
		fmt.Println("Please provide a remote origin as argument.")
		os.Exit(1)
	}

	remoteOrigin := args[0]

	// Check if Dotfiles folder exists
	_, err := os.Stat(repoPath)
	if errors.Is(err, fs.ErrNotExist) {
		// Create directory if it does not exist
		err = os.Mkdir(repoPath, fs.ModePerm)
		if err != nil {
			fmt.Println("Unable creating directory:", err)
			os.Exit(1)
		}

		fmt.Printf("Created directory %s\n", repoPath)
	} else {
		fmt.Printf("Directory %s already exists\n", repoPath)
	}

	// Check if .badm.yaml file exists in .dotfiles folder
	_, err = os.Stat(cfgFile)
	if errors.Is(err, fs.ErrNotExist) {
		// Create .badm.yaml if it does not exist
		file, err := os.Create(cfgFile)
		if err != nil {
			fmt.Println("Unable creating .badm.yaml config file:", err)
			os.Exit(1)
		}

		defer file.Close()

		fmt.Println("Created .badm.yaml")
	} else {
		fmt.Println("Config file .badm.yaml already exists")
	}

	// Check if .gitignore file exists in .dotfiles folder
	gitignore := fmt.Sprintf("%s/%s", repoPath, ".gitignore")

	_, err = os.Stat(gitignore)
	if errors.Is(err, fs.ErrNotExist) {
		// Create .gitignore if it does not exist
		err = os.WriteFile(gitignore, []byte("badm.state"), 0600) //nolint:mnd
		if err != nil {
			fmt.Println("Unable creating .gitignore config file:", err)
			os.Exit(1) //nolint:gocritic
		}

		fmt.Println("Created .gitignore")
	} else {
		fmt.Println("Config file .gitignore already exists")
	}

	git.InitGitRepository(repoPath, remoteOrigin)

	fmt.Printf("Set up Git repository with remote '%s'\n", remoteOrigin)

	fmt.Println("Pushing config files to remote...")

	git.CommitAndPushFiles(repoPath, "🚀 Add initial config files")

	fmt.Println("Pushed config files to remote")
}
