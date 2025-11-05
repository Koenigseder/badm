package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CommitAndPushFiles commits and pushes all untracked files in the BADM Git repository
func CommitAndPushFiles(repoPath, commitMsg string) {
	filesChanged, err := checkIfFilesChanged(repoPath)
	if err != nil {
		fmt.Println("Unable checking file status:", err)
		os.Exit(1)
	}

	if !filesChanged {
		fmt.Println("No files changed! Everything's up to date :)")
		os.Exit(0)
	}

	// git add
	_, err = exec.Command("git", "-C", repoPath, "add", ".").Output()
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

func checkIfFilesChanged(repoPath string) (bool, error) {
	output, err := exec.Command("git", "-C", repoPath, "status").Output()
	if err != nil {
		return false, err
	}

	if strings.Contains(string(output), "nothing to commit") {
		return false, nil
	}

	return true, nil
}
