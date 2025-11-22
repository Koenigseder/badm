package filesystem

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

// MoveFileWithSymLink moves a file from [sourcePath] to [destinationPath] and creates a soft symbolic link at [sourcePath]
func MoveFileWithSymLink(sourcePath, destinationPath string, withSudo bool) {
	dirPath, _ := path.Split(destinationPath)

	mkdirCmd := exec.Command("mkdir", "-p", dirPath)
	mvCmd := exec.Command("mv", sourcePath, destinationPath)
	lnCmd := exec.Command("ln", "-s", destinationPath, sourcePath)

	if withSudo {
		mkdirCmd = exec.Command("sudo", "mkdir", "-p", dirPath)
		mvCmd = exec.Command("sudo", "mv", sourcePath, destinationPath)
		lnCmd = exec.Command("sudo", "ln", "-s", destinationPath, sourcePath)
	}

	// Create directories if necessary
	err := mkdirCmd.Run()
	if err != nil {
		fmt.Println("Unable creating needed directories:", err)
		os.Exit(1)
	}

	// Move file
	err = mvCmd.Run()
	if err != nil {
		fmt.Println("Unable moving file:", err)
		os.Exit(1)
	}

	// Create soft symbolic link
	err = lnCmd.Run()
	if err != nil {
		fmt.Println("Unable creating symbolic link:", err)
		os.Exit(1)
	}
}
