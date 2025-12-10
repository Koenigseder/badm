package cmd

import (
	"fmt"
	"os"

	pkgs "github.com/Koenigseder/badm/internal/packages"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(packages)
}

var packages = &cobra.Command{
	Use:   "packages",
	Short: "Install various packages",
	Long:  `Install various packages defined in .badm.yaml`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a install scope name")
			os.Exit(1)
		}

		installScopeName := args[0]

		cfgFile, err := pkgs.ReadConfigFile(cfgFilePath)
		if err != nil {
			fmt.Printf("Failed reading config file %s: %v\n", cfgFilePath, err)
			os.Exit(1)
		}

		cfgFile.Directory = repoPath

		err = cfgFile.InstallPackages(installScopeName)
		if err != nil {
			fmt.Printf("Failed installing packages for '%s': %v\n", installScopeName, err)
			os.Exit(1)
		}
	},
}
