package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// These variables are filled via -ldflags on compile time
var version string
var revision string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `All software has versions. This is BADM's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("BADM - Born Again Dotfile Manager\n Version: %s - Revision: %s\n", version, revision)
	},
}
