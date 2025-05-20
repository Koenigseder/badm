package filesystem

import (
	"encoding/json"
	"fmt"
	"os"
)

// The State object
type State struct {
	Name       string   `json:"name"`
	LocalFiles []string `json:"local_files"`
}

// WriteStateFile writes a State struct object to disc
func (s *State) WriteStateFile(repoBasePath string) error {
	jsonString, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	stateFilePath := fmt.Sprintf("%s/%s", repoBasePath, s.Name)

	err = os.WriteFile(stateFilePath, jsonString, 0600) //nolint:mnd
	if err != nil {
		return err
	}

	return nil
}

// ReadStateFile reads a state file's content into a pre-allocated State struct
func (s *State) ReadStateFile(repoBasePath string) error {
	if s.Name == "" {
		s.Name = "badm.state"
	}

	stateFilePath := fmt.Sprintf("%s/%s", repoBasePath, s.Name)

	fileBytes, err := os.ReadFile(stateFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileBytes, &s)
	if err != nil {
		return err
	}

	return nil
}
