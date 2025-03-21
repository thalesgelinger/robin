package entities

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Build struct {
		IOS IOS `yaml:"ios"`
	} `yaml:"build"`
}

type IOS struct {
	AppName     string                 `yaml:"app_name"`
	ProjectPath string                 `yaml:"project_path"`
	OutputDir   string                 `yaml:"output_dir"`
	Envs        map[string]Environment `yaml:"environments"`
}

type Environment struct {
	Scheme               string `yaml:"scheme"`
	IncrementBuildNumber bool   `yaml:"increment_build_number"`
	ExportMethod         string `yaml:"export_method"`
}

func ReadConfig() Config {
	file, err := os.ReadFile("robin.yml")
	if err != nil {
		fmt.Println("Error reading robin.yml:", err)
		os.Exit(1)
	}

	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		fmt.Println("Error parsing YAML:", err)
		os.Exit(1)
	}
	return config
}

