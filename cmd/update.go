package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Use the update command to update previously installed applications.",
	Long: `The update command opens an MQTT connection to the
server, receives update command on new commits to Github, and
orchestrates updating the Systemd unit.`,
	Run: func(cmd *cobra.Command, args []string) {
		runUpdate(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) {
	cfg := getConfig(cmd)
	herokuAPIKey := os.Getenv("HEROKU_API_KEY")
	if herokuAPIKey == "" {
		logger.Fatalln("HEROKU_API_TOKEN environment variable is required")
	}

	_, err := newAgent(herokuAPIKey)
	if err != nil {
		logger.Fatalln(fmt.Errorf("error creating agent: %s", err))
	}

	fmt.Println(cfg)
}
