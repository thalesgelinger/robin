package entities

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Build `yaml:"build"`
}

type Build struct {
	IOS `yaml:"ios"`
}

type IOS struct {
	AppName string            `yaml:"app_name"`
	Envs    map[string]IosEnv `yaml:"environments"`
}

type IosEnv struct {
	IncrementBuildNumber bool   `yaml:"increment_build_number"`
	Scheme               string `yaml:"scheme"`
	ExportMethod         string `yaml:"export_method"`
}

func ReadConfig() Config {

	file, err := os.ReadFile("robin.yml")
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		fmt.Printf("Error parsing YAML: %v\n", err)
		os.Exit(1)
	}

	return config
}
