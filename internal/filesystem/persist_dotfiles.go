package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// PersistDotfiles reads all Dotfiles and creates symlinks for everyone on the system
// returns a string slice with all persisted files
func PersistDotfiles(baseDir, repoName, repoRootAlias string, overrideExistingFiles bool) []string {
	writtenFiles := new([]string)

	err := filepath.WalkDir(baseDir, func(dotfilesPath string, d fs.DirEntry, _ error) error {
		// Skip directories, .git folder, badm.state and .badm.yaml
		if d.IsDir() ||
			(strings.Contains(dotfilesPath, ".git") && !strings.HasSuffix(dotfilesPath, ".gitconfig")) ||
			strings.HasSuffix(dotfilesPath, ".badm.yaml") ||
			strings.HasSuffix(dotfilesPath, "badm.state") {
			return nil
		}

		var (
			destinationPath string
			mkdirCmd        *exec.Cmd
			rmCmd           *exec.Cmd
			lnCmd           *exec.Cmd
		)

		if strings.Contains(dotfilesPath, repoRootAlias) {
			// Remove Dotfiles repo path and `repoRootAlias` from path
			destinationPath = filepath.Clean(strings.Replace(dotfilesPath, baseDir+"/"+repoRootAlias, "", 1))
			dirPath, _ := path.Split(destinationPath)

			mkdirCmd = exec.Command("sudo", "mkdir", "-p", dirPath)
			rmCmd = exec.Command("sudo", "rm", destinationPath)
			lnCmd = exec.Command("sudo", "ln", "-s", dotfilesPath, destinationPath)
		} else {
			// Remove Dotfiles repo name from path
			destinationPath = filepath.Clean(strings.Replace(dotfilesPath, repoName, "", 1))
			dirPath, _ := path.Split(destinationPath)

			mkdirCmd = exec.Command("mkdir", "-p", dirPath)
			rmCmd = exec.Command("rm", destinationPath)
			lnCmd = exec.Command("ln", "-s", dotfilesPath, destinationPath)
		}

		// Create directories if necessary
		err := mkdirCmd.Run()
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
			err = rmCmd.Run()
			if err != nil {
				fmt.Printf("Unable deleting existing file %s: %v\n", destinationPath, err)
				os.Exit(1)
			}

			fmt.Println("Removed file", destinationPath)
		}

		// Create soft symbolic link
		err = lnCmd.Run()
		if err != nil {
			fmt.Println("Unable creating symbolic link:", err)
			os.Exit(1)
		}

		fmt.Println("Created symlink:", destinationPath)

		// Append the written file to string slice
		*writtenFiles = append(*writtenFiles, destinationPath)

		return nil
	})
	if err != nil {
		fmt.Println("Unable walking through Dotfile repo", err)
		os.Exit(1)
	}

	return *writtenFiles
}
