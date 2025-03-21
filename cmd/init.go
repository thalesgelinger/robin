package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Robin configuration",
	Long:  `Initialize Robin by generating a default robin.yml based on the current project type.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectType := detectProjectType()

		switch projectType {
		case "expo":
			fmt.Println("📦 Detected Expo project")
			initializeExpoProject()
		default:
			fmt.Println("⚠️ Could not detect a known project type. Creating a default robin.yml...")
			createDefaultConfig()
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// Detect project type (currently supports Expo, ready for expansion)
func detectProjectType() string {
	if _, err := os.Stat("app.json"); err == nil {
		return "expo"
	}
	return "unknown"
}

// Initialize an Expo project by extracting info from app.json
func initializeExpoProject() {
	expoConfig := readExpoConfig()
	if expoConfig == nil {
		fmt.Println("❌ Failed to read app.json")
		return
	}

	appName := expoConfig["app_name"]
	scheme := expoConfig["scheme"]
	if scheme == "" {
		scheme = appName // Default scheme = app name (lowercased)
	}

	// Create robin.yml
	config := map[string]any{
		"build": map[string]any{
			"ios": map[string]any{
				"app_name":     appName,
				"project_path": "./ios",
				"output_dir":   "./builds",
				"environments": map[string]any{
					"development": map[string]any{
						"scheme":                 scheme,
						"increment_build_number": true,
						"export_method":          "development",
					},
					"production": map[string]any{
						"scheme":                 scheme,
						"increment_build_number": true,
						"export_method":          "app-store",
					},
				},
			},
		},
	}

	writeConfigFile(config)
	fmt.Println("✅ robin.yml created successfully!")
}

// Read app.json for Expo projects
func readExpoConfig() map[string]string {
	data, err := os.ReadFile("app.json")
	if err != nil {
		fmt.Println("❌ Error reading app.json:", err)
		return nil
	}

	var appJson map[string]any
	if err := json.Unmarshal(data, &appJson); err != nil {
		fmt.Println("❌ Error parsing app.json:", err)
		return nil
	}

	expoConfig := make(map[string]string)

	// Extract "expo" config
	if expo, ok := appJson["expo"].(map[string]any); ok {
		// Get app name (first try "name", then fallback to "slug")
		if name, found := expo["name"].(string); found {
			expoConfig["app_name"] = name
		} else if slug, found := expo["slug"].(string); found {
			expoConfig["app_name"] = slug
		}

		// Get scheme
		if scheme, found := expo["scheme"].(string); found {
			expoConfig["scheme"] = scheme
		}
	}

	return expoConfig
}

// Create a default robin.yml file for unknown projects
func createDefaultConfig() {
	config := map[string]any{
		"build": map[string]any{
			"ios": map[string]any{
				"app_name":     "MyApp",
				"project_path": "./ios",
				"output_dir":   "./builds",
				"environments": map[string]any{
					"development": map[string]any{
						"scheme":                 "MyApp",
						"increment_build_number": true,
						"export_method":          "development",
					},
					"production": map[string]any{
						"scheme":                 "MyApp",
						"increment_build_number": true,
						"export_method":          "app-store",
					},
				},
			},
		},
	}

	writeConfigFile(config)
	fmt.Println("✅ Default robin.yml created successfully!")
}

// Write robin.yml to file
func writeConfigFile(config map[string]any) {
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		fmt.Println("❌ Failed to generate YAML:", err)
		return
	}

	err = os.WriteFile("robin.yml", yamlData, 0644)
	if err != nil {
		fmt.Println("❌ Error writing robin.yml:", err)
		return
	}
}

