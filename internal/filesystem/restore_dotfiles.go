package filesystem

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// RestoreDotfiles copies all original Dotfiles to the system and removes symlinks
func RestoreDotfiles(baseDir, repoName string, dryRun bool) {
	err := filepath.WalkDir(baseDir, func(dotfilesPath string, d fs.DirEntry, err error) error {
		// Skip directories, .git folder and .badm.yaml
		if d.IsDir() || strings.Contains(dotfilesPath, ".git") || strings.HasSuffix(dotfilesPath, ".badm.yaml") {
			return nil
		}

		// Remove Dotfiles repo name from path
		destinationPath := filepath.Clean(strings.Replace(dotfilesPath, repoName, "", 1))

		// Check if file exists
		if _, err = os.Stat(destinationPath); err == nil {
			// Check if file should not be restored but output (Dry Run)
			if dryRun {
				fmt.Println("Restored file", destinationPath, "(DRY)")

				return nil
			}

			// Remove symlink
			err = os.Remove(destinationPath)
			if err != nil {
				fmt.Printf("Unable removing symlink %s: %v\n", destinationPath, err)
				os.Exit(1)
			}

			// Copy file
			srcFile, err := os.Open(dotfilesPath)
			if err != nil {
				fmt.Printf("Unable opening file %s: %v\n", dotfilesPath, err)
				os.Exit(1)
			}
			defer srcFile.Close()

			destFile, err := os.Create(destinationPath)
			if err != nil {
				fmt.Printf("Unable creating file %s: %v\n", destinationPath, err)
				os.Exit(1)
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, srcFile)
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
