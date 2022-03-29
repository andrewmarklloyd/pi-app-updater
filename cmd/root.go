package cmd

import (
	"log"
	"os"
	"os/user"

	"github.com/spf13/cobra"
)

var logger = log.New(os.Stdout, "[pi-app-deployer-Agent] ", log.LstdFlags)

var rootCmd = &cobra.Command{
	Use:   "pi-app-deployer-agent",
	Short: "",
	Long:  ``,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	u, err := user.Current()
	if err != nil {
		logger.Fatalln("error getting current user:", err)
	}
	if u.Username != "root" {
		logger.Fatalln("agent must be run as root, user found was", u.Username)
	}

	rootCmd.PersistentFlags().String("repoName", "", "Name of the Github repo including the owner")
	rootCmd.PersistentFlags().String("manifestName", "", "Name of the pi-app-deployer manifest")
	rootCmd.PersistentFlags().Bool("logForwarding", false, "Send application logs to server")
	rootCmd.PersistentFlags().String("appUser", "pi", "Name of user that will run the app service")

	rootCmd.PersistentFlags().Var(&varFlags, "envVar", "List of non-secret environment variable configuration, separated by =, can pass multiple values. Example: --env-var foo=bar --env-var hello=world")
}
