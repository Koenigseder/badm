package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// These variables are filled via -ldflags on compile time
var (
	version  string
	revision string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `All software has versions. This is BADM's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("BADM - Born Again Dotfile Manager\n Version: %s - Revision: %s\n", version, revision)
	},
}

func checkRemoteVersion() {
	type latestReleaseResponse struct {
		TagName string `json:"tag_name"`
	}

	badmLatestReleaseUrl := "https://api.github.com/repos/Koenigseder/badm/releases/latest"

	res, err := http.Get(badmLatestReleaseUrl)
	if err != nil {
		fmt.Println("Unable getting latest release:", err)
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Unable reading release API response:", err)
		return
	}

	latestReleaseRes := new(latestReleaseResponse)

	err = json.Unmarshal(resBody, &latestReleaseRes)
	if err != nil {
		fmt.Println("Unable parsing release API response to JSON:", err)
		return
	}

	if semver.Compare(version, latestReleaseRes.TagName) == -1 {
		formatString(latestReleaseRes.TagName)
	}
}

func formatString(latestReleaseTag string) {
	newVersionText := fmt.Sprintf("A new version is available! (%s)", latestReleaseTag)
	downloadText := fmt.Sprintf("     Download at https://github.com/Koenigseder/badm/releases/tag/%s", latestReleaseTag)

	sharpLength := len(downloadText) + 2*3 // Length of header and trailer

	fmt.Println("\033[33m")

	for j := 0; j < sharpLength; j++ {
		fmt.Printf("#")
	}
	fmt.Println()

	for i := 0; i < len(newVersionText)/2; i++ {
		fmt.Printf(" ")
	}

	fmt.Println(newVersionText)
	fmt.Println(downloadText)

	for j := 0; j < sharpLength; j++ {
		fmt.Printf("#")
	}
	fmt.Println("\n\033[0m")
}
