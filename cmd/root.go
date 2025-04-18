package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	repoName = ".dotfiles_badm"

	var err error

	homeDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Println("Error retrieving user's home directory:", err)
		os.Exit(1)
	}

	repoPath = fmt.Sprintf("%s/%s", homeDir, repoName)
	cfgFile = fmt.Sprintf("%s/.badm.yaml", repoPath)
}

var (
	repoName string
	homeDir  string
	repoPath string
	cfgFile  string

	// Flags
	overrideExistingFiles bool

	rootCmd = &cobra.Command{
		Use:   "badm",
		Short: "BADM is a Dotfile manager",
		Long:  `BADM - Born Again Dotfile Manager`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("No argument provided!")
		},
	}
)

// Execute the app
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
