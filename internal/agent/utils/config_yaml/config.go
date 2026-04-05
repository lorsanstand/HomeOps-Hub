package config_yaml

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AgentConfig struct {
	AppName    string `yaml:"app_name"`
	HubConnect struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"hub"`
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
