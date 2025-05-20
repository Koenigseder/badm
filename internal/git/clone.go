package git

import (
	"fmt"
	"os"
	"os/exec"
)

// CloneGitRepository clones the Git repository
func CloneGitRepository(repoURL, repoPath string) {
	cmd := exec.Command("git", "clone", repoURL, repoPath)

	fmt.Printf("Getting %s...\n", repoURL)

	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("Unable cloning Git repository %s: %v\n", repoURL, err)
		os.Exit(1)
	}

	fmt.Printf("Got %s\n", repoURL)
}
