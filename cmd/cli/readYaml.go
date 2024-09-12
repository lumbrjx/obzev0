package main

import (
	"fmt"
	"obzev0/common/definitions"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfig(path string) (definitions.Config, error) {
	var config definitions.Config

	yamlFile, err := os.ReadFile(path + "obzev0cnf.yaml")
	if err != nil {
		return config, fmt.Errorf("failed to read YAML file: %w", err)
	}

	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return config, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return config, nil
}
