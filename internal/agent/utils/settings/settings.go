package settings

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type Settings struct {
	AgentID string `json:"agent_id"`
	path    string
}

func ReadSettings(path string) (*Settings, error) {
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(homeDir, ".config", "homeops")
	}

	err := os.Mkdir(path, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		err = nil
	}

	file, err := os.Create(path + "/settings.json")
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		err = nil
	}
	defer file.Close()

	var settings Settings

	err = json.NewDecoder(file).Decode(&settings)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}

	settings.path = path + "/settings.json"

	return &settings, nil
}

func (s *Settings) Insert(sett Settings) error {
	file, err := os.OpenFile(s.path, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = json.NewEncoder(file).Encode(sett); err != nil {
		return err
	}

	return nil
}
