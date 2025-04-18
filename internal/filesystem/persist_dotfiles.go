package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// PersistDotfiles reads all Dotfiles and creates symlinks for everyone on the system
func PersistDotfiles(baseDir, repoName string, overrideExistingFiles bool) {
	err := filepath.WalkDir(baseDir, func(dotfilesPath string, d fs.DirEntry, err error) error {
		// Skip directories, .git folder and .badm.yaml
		if d.IsDir() || strings.Contains(dotfilesPath, ".git") || strings.HasSuffix(dotfilesPath, ".badm.yaml") {
			return nil
		}

		// Remove Dotfiles repo name from path
		destinationPath := filepath.Clean(strings.Replace(dotfilesPath, repoName, "", 1))

		// Create directories if necessary
		dirPath, _ := path.Split(destinationPath)
		err = os.MkdirAll(dirPath, fs.ModePerm)
		if err != nil {
			fmt.Println("Unable creating needed directories:", err)
			os.Exit(1)
		}

		// Check if file exists
		if _, err = os.Stat(destinationPath); err == nil {
			// Abort if file should not be overwritten
			if !overrideExistingFiles {
				fmt.Println("File exists:", destinationPath)

				return nil
			}

			// Override file
			err = os.Remove(destinationPath)
			if err != nil {
				fmt.Printf("Unable deleting existing file %s: %v\n", destinationPath, err)
				os.Exit(1)
			}

			fmt.Println("Removed file", destinationPath)
		}

		// Create soft symbolic link
		err = os.Symlink(dotfilesPath, destinationPath)
		if err != nil {
			fmt.Println("Unable creating symbolic link:", err)
			os.Exit(1)
		}

		fmt.Println("Created symlink:", destinationPath)

		return nil
	})
	if err != nil {
		fmt.Println("Unable walking through Dotfile repo", err)
		os.Exit(1)
	}
}
