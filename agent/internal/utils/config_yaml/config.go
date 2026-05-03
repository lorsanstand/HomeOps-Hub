package config_yaml

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

const MODE = "DEV"
const AgentVersion = "0.0"

type AgentConfig struct {
	AppName    string `yaml:"app_name"`
	HubConnect struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"hub"`
	LogLevel     string `yaml:"log_level"`
	SettingsPath string `yaml:"settings_path"`
}

func NewConfig() (*AgentConfig, error) {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed open file: %v", err)
	}

	var cfg AgentConfig

	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, fmt.Errorf("failed read yaml: %v", err)
	}

	return &cfg, nil
}

func (c *AgentConfig) GetLogLevel() zerolog.Level {
	level, err := zerolog.ParseLevel(c.LogLevel)
	if err != nil {
		return zerolog.InfoLevel
	}
	return level
}

func (c *AgentConfig) GetMode() string {
	return MODE
}

func (c *AgentConfig) GetAgentVersion() string {
	return AgentVersion
}

func (c *AgentConfig) GetGRPCAddress() string {
	return fmt.Sprintf("%v:%v", c.HubConnect.Host, c.HubConnect.Port)
}
