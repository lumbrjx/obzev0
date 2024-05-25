package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Delays DelaysConfig `yaml:"delays"`
	Server ServerConfig `yaml:"server"`
}

type DelaysConfig struct {
	ReqDelay int `yaml:"reqDelay"`
	ResDelay int `yaml:"resDelay"`
}
type ServerConfig struct {
	Port string `yaml:"port"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config

	// Read YAML file
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Unmarshal YAML into struct
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return config, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return config, nil
}
