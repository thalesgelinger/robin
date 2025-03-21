package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate ",
	Long:  `Logs into using Apple ID and fetches the App Store Connect team ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		// email := promptInput("Enter Apple ID Email")
		// password := promptInput("Enter Apple ID Password")
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}

func promptInput(prompt string) string {
	fmt.Print(prompt + ": ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

