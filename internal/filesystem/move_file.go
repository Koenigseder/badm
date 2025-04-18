package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path"
)

// MoveFileWithSymLink moves a file from [sourcePath] to [destinationPath] and creates a soft symbolic link at [sourcePath]
func MoveFileWithSymLink(sourcePath, destinationPath string) {
	// Create directories if necessary
	dirPath, _ := path.Split(destinationPath)
	err := os.MkdirAll(dirPath, fs.ModePerm)
	if err != nil {
		fmt.Println("Unable creating needed directories:", err)
		os.Exit(1)
	}

	// Move file
	err = os.Rename(sourcePath, destinationPath)
	if err != nil {
		fmt.Println("Unable moving file:", err)
		os.Exit(1)
	}

	// Create soft symbolic link
	err = os.Symlink(destinationPath, sourcePath)
	if err != nil {
		fmt.Println("Unable creating symbolic link:", err)
		os.Exit(1)
	}
}
