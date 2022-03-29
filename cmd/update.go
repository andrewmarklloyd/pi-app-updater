package cmd

import (
	"fmt"

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
		fmt.Println("update called")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
