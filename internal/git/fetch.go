package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// FetchAndUpdate fetches and pulls the latest changes from the remote repository.
// Returns "true" if there are any changes
func FetchAndUpdate(repoPath string) bool {
	// git fetch
	_, err := exec.Command("git", "-C", repoPath, "fetch").Output()
	if err != nil {
		fmt.Println("Unable fetching Git repository:", err)
		os.Exit(1)
	}

	// git pull
	out, err := exec.Command("git", "-C", repoPath, "pull").Output()
	if err != nil {
		fmt.Println("Unable pulling Git repository:", err)
		os.Exit(1)
	}

	if !strings.HasPrefix(string(out), "Already up to date.") {
		return true
	}

	return false
}
