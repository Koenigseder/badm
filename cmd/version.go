package cmd

import (
	"fmt"
	"strings"

	"github.com/Koenigseder/badm/internal/core"
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
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("BADM - Born Again Dotfile Manager\n Version: %s - Revision: %s\n", version, revision)
	},
}

func checkRemoteVersion() {
	latestReleaseInformation, err := core.GetLatestReleaseInformation()
	if err != nil {
		fmt.Println("Unable getting latest release information:", err)
		return
	}

	if semver.Compare(version, latestReleaseInformation.TagName) == -1 {
		formatString(latestReleaseInformation.TagName)
	}
}

func formatString(latestReleaseTag string) {
	newVersionText := fmt.Sprintf("A new version is available! (%s)", latestReleaseTag)
	downloadText := fmt.Sprintf("   Update with \033[1m\033[32mbadm update\033[0m\033[33m or download at \033[1m\033[32mhttps://github.com/Koenigseder/badm/releases/tag/%s\033[0m\033[33m", latestReleaseTag)

	visibleNewVersionLength := visibleLength(newVersionText)
	visibleDownloadLength := visibleLength(downloadText)

	// Length of header and trailer
	frameLength := visibleDownloadLength + 3 //nolint:mnd

	frame := strings.Repeat("#", frameLength)

	centeredNewVersion := fmt.Sprintf("%[1]*s", -frameLength, fmt.Sprintf("%[1]*s", (frameLength+visibleNewVersionLength)/2, newVersionText))
	centeredDownload := fmt.Sprintf("%[1]*s", -frameLength, fmt.Sprintf("%[1]*s", (frameLength+visibleDownloadLength)/2, downloadText))

	fmt.Println("\033[33m" + frame)
	fmt.Println(centeredNewVersion)
	fmt.Println(centeredDownload)
	fmt.Println(frame + "\033[0m")
}

func visibleLength(s string) int {
	// ANSI-Escape-Sequences start with '\033[' oder '\x1b['
	// We only count the visible characters
	length := 0
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\033' || s[i] == '\x1b' {
			inEscape = true
			continue
		}

		if inEscape {
			if s[i] >= 'A' && s[i] <= 'Z' || s[i] == 'm' {
				inEscape = false
			}

			continue
		}

		length++
	}

	return length
}
