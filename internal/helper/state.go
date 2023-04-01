package helper

import (
	"encoding/json"
	"os"
)

type State struct {
	SSHPrivateKey string
	ServerAddress string
}

const StatePath = "state.json"

func (s *State) WriteToFile(path string) error {
	os.Truncate(path, 0)
	fh, fileOpenError := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	if fileOpenError != nil {
		return fileOpenError
	}

	return json.NewEncoder(fh).Encode(s)
}

func ReadStateFromFile(path string) (*State, error) {
	fh, fileOpenError := os.Open(path)
	if fileOpenError != nil {
		return nil, fileOpenError
	}

	var state State
	if err := json.NewDecoder(fh).Decode(&state); err != nil {
		return nil, err
	}

	return &state, nil
}
