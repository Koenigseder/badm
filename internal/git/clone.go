package git

import (
	"fmt"
	"os"
	"os/exec"
)

// CloneGitRepository clones the Git repository
func CloneGitRepository(repoUrl, repoPath string) {
	cmd := exec.Command("git", "clone", repoUrl, repoPath)

	fmt.Printf("Getting %s...\n", repoUrl)
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("Unable cloning Git repository %s: %v\n", repoUrl, err)
		os.Exit(1)
	}

	fmt.Printf("Got %s\n", repoUrl)
}
