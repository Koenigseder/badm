package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ReleaseInformation partially represents the API response of a release information fetch
type ReleaseInformation struct {
	TagName string `json:"tag_name"`
}

// GetLatestReleaseInformation fetches information about the latest BADM release and returns it as a struct
func GetLatestReleaseInformation() (*ReleaseInformation, error) {
	latestReleaseInformation := new(ReleaseInformation)

	badmLatestReleaseInformationURL := "https://api.github.com/repos/Koenigseder/badm/releases/latest"

	res, err := http.Get(badmLatestReleaseInformationURL)
	if err != nil {
		fmt.Println("Unable getting latest release:", err)
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Unable reading release API response:", err)
		return nil, err
	}

	err = json.Unmarshal(resBody, &latestReleaseInformation)
	if err != nil {
		fmt.Println("Unable parsing release API response to JSON:", err)
		return nil, err
	}

	return latestReleaseInformation, nil
}
