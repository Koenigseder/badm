package git

import (
	"fmt"
	"os"
	"os/exec"
)

// InitGitRepository initializes new Git repository
func InitGitRepository(repoPath, remoteOrigin string) {
	// git init
	_, err := exec.Command("git", "-C", repoPath, "init").Output()
	if err != nil {
		fmt.Println("Unable initializing Git repository:", err)
		os.Exit(1)
	}

	// git branch
	_, err = exec.Command("git", "-C", repoPath, "branch", "-M", "linux").Output()
	if err != nil {
		fmt.Println("Unable setting Git branch:", err)
		os.Exit(1)
	}

	// git remote add
	_, err = exec.Command("git", "-C", repoPath, "remote", "add", "origin", remoteOrigin).Output()
	if err != nil {
		fmt.Println("Unable adding Git remote:", err)
		os.Exit(1)
	}
}
