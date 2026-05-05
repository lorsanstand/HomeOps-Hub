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

func ReadSettings(path string) (sett *Settings, err error) {
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(homeDir, ".config", "homeops")
	}

	err = os.Mkdir(path, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		err = nil
	}

	settingsPath := filepath.Join(path, "settings.json")
	var settings Settings

	file, err := os.Open(settingsPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	} else {
		defer func() {
			closeErr := file.Close()
			if err == nil {
				err = closeErr
			}
		}()
		err = json.NewDecoder(file).Decode(&settings)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
	}

	settings.path = settingsPath

	return &settings, nil
}

func (s *Settings) InsertAgentID(agentID string) (err error) {
	file, err := os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := file.Close()
		if err == nil {
			err = closeErr
		}
	}()

	sett := Settings{AgentID: agentID}

	if err = json.NewEncoder(file).Encode(sett); err != nil {
		return err
	}

	return nil
}

func (s *Settings) GetAgentID() string {
	return s.AgentID
}
