package services

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Credentials struct {
	AppIdentifier string `yaml:"app_identifier"`
	AppleID       string `yaml:"apple_id"`
	ItcTeamId     string `yaml:"itc_team_id"`
	TeamID        string `yaml:"team_id"`
}

func LoadCredentials() *Credentials {
	file, err := os.ReadFile("robin_credentials.yml")
	if err != nil {
		fmt.Println("❌ Error reading robin_credentials.yml:", err)
		os.Exit(1)
	}

	var config Credentials
	if err := yaml.Unmarshal(file, &config); err != nil {
		fmt.Println("❌ Error parsing YAML:", err)
		os.Exit(1)
	}
	return &config
}
