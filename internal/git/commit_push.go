package git

import (
	"fmt"
	"os"
	"os/exec"
)

// CommitAndPushFiles commits and pushes all untracked files in the BADM Git repository
func CommitAndPushFiles(repoPath, commitMsg string) {
	// git add
	_, err := exec.Command("git", "-C", repoPath, "add", ".").Output()
	if err != nil {
		fmt.Println("Unable adding Dotfile to Git:", err)
		os.Exit(1)
	}

	// git commit
	_, err = exec.Command("git", "-C", repoPath, "commit", "-m", commitMsg).Output()
	if err != nil {
		fmt.Println("Unable committing Dotfile to Git:", err)
		os.Exit(1)
	}

	// git push
	_, err = exec.Command("git", "-C", repoPath, "push", "-u", "origin", "linux").Output()
	if err != nil {
		fmt.Println("Unable pushing Dotfile to remote Git repository:", err)
		os.Exit(1)
	}
}
