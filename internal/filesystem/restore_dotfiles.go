package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RestoreDotfiles copies all original Dotfiles to the system and removes symlinks
func RestoreDotfiles(baseDir, repoName, repoRootAlias string, dryRun bool) {
	err := filepath.WalkDir(baseDir, func(dotfilesPath string, d fs.DirEntry, _ error) error {
		// Skip directories, .git folder, .badm.yaml and badm.state
		if d.IsDir() ||
			(strings.Contains(dotfilesPath, ".git") && !strings.HasSuffix(dotfilesPath, ".gitconfig")) ||
			strings.HasSuffix(dotfilesPath, ".badm.yaml") ||
			strings.HasSuffix(dotfilesPath, "badm.state") {
			return nil
		}

		var (
			destinationPath string
			rmCmd           *exec.Cmd
			cpCmd           *exec.Cmd
		)

		// Remove Dotfiles repo name from path
		if strings.Contains(dotfilesPath, repoRootAlias) {
			stringToRemove := baseDir + "/" + repoRootAlias
			destinationPath = filepath.Clean(strings.Replace(dotfilesPath, stringToRemove, "", 1))

			rmCmd = exec.Command("sudo", "rm", destinationPath)
			cpCmd = exec.Command("sudo", "cp", dotfilesPath, destinationPath)
		} else {
			destinationPath = filepath.Clean(strings.Replace(dotfilesPath, repoName, "", 1))

			rmCmd = exec.Command("rm", destinationPath)
			cpCmd = exec.Command("cp", dotfilesPath, destinationPath)
		}

		// Check if file exists
		if _, err := os.Stat(destinationPath); err == nil {
			// Check if file should not be restored but output (Dry Run)
			if dryRun {
				fmt.Println("Restored file", destinationPath, "(DRY)")

				return nil
			}

			// Remove symlink
			err = rmCmd.Run()
			if err != nil {
				fmt.Printf("Unable removing symlink %s: %v\n", destinationPath, err)
				os.Exit(1)
			}

			// Copy file
			err = cpCmd.Run()
			if err != nil {
				fmt.Printf("Unable copying file %s: %v\n", dotfilesPath, err)
				os.Exit(1)
			}

			fmt.Println("Restored file", destinationPath)

			return nil
		}

		fmt.Println("File not found:", destinationPath)

		return nil
	})
	if err != nil {
		fmt.Println("Unable walking through Dotfile repo", err)
		os.Exit(1)
	}
}
