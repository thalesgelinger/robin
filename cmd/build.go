package cmd

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/thalesgelinger/robin/internal/entities"
	"github.com/thalesgelinger/robin/internal/services"
)

var verbose bool

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build app",
	Long:  `Build app using the configuration from robin.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		config := entities.ReadConfig()

		if len(args) == 0 {
			fmt.Println("Error: You must specify a platform (ios).")
			return
		}

		platform := args[0]

		fmt.Printf("Building app: %s\n", config.Build.IOS.AppName)

		switch platform {
		case "ios":
			runIosBuild(config.Build.IOS, verbose)
		default:
			fmt.Println("Unsupported platform:", platform)
		}
	},
}

func init() {
	buildCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.AddCommand(buildCmd)
}

func runIosBuild(ios entities.IOS, verbose bool) {
	env := ios.Envs["development"]

	fmt.Println("🔧 Using Scheme:", env.Scheme)

	if env.IncrementBuildNumber {
		incrementBuildNumber(ios.ProjectPath, env.Scheme)
	}

	credentials := services.LoadCredentials()

	xcode, err := services.NewXcode(ios, env, credentials, verbose)

	if err != nil {
		log.Fatalf("❌ Error generating ExportPlist: %v", err)
	}

	if err := xcode.Archive(); err != nil {
		log.Fatalf("❌ Archive failed: %v", err)
	}

	if err := xcode.ExportIPA(); err != nil {
		log.Fatalf("❌ IPA export failed: %v", err)
	}

	fmt.Println("✅ Build completed successfully!")
}

func incrementBuildNumber(projectPath, scheme string) {
	fmt.Println("🔼 Incrementing build number...")
	cmd := exec.Command("agvtool", "next-version", "-all")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Failed to increment build number:", err)
	} else {
		fmt.Println("✅ Build number incremented:", string(output))
	}
}
