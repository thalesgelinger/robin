package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thalesgelinger/robin/cmd/entities"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build app",
	Long:  `Build app`,
	Run: func(cmd *cobra.Command, args []string) {
		config := entities.ReadConfig()

		if len(args) == 0 {
			fmt.Printf("You should inform a platform to run build")
			return
		}

		platform := args[0]

		fmt.Printf("Building app: %s\n", config.Build.IOS.AppName)

		switch platform {
		case "ios":
			runIosBuild(config.Build.IOS)
		}

	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func runIosBuild(ios entities.IOS) {
	env := ios.Envs["development"]
	fmt.Println("Scheme ", env.Scheme)
}

