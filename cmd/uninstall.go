package cmd

import (
	"fmt"
	"os"

	"github.com/andrewmarklloyd/pi-app-deployer/internal/pkg/config"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "TODO",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		runUninstall(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	uninstallCmd.PersistentFlags().Bool("all", false, "Uninstall all apps")
	uninstallCmd.PersistentFlags().String("repoName", "", "Name of the Github repo including the owner")
	uninstallCmd.PersistentFlags().String("manifestName", "", "Name of the pi-app-deployer manifest")
}

func runUninstall(cmd *cobra.Command, args []string) {
	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		fmt.Println("error getting all flag", err)
		os.Exit(1)
	}

	repoName, err := cmd.Flags().GetString("repoName")
	if err != nil {
		fmt.Println("error getting repoName flag", err)
		os.Exit(1)
	}

	manifestName, err := cmd.Flags().GetString("manifestName")
	if err != nil {
		fmt.Println("error getting manifestName flag", err)
		os.Exit(1)
	}

	herokuApp, err := cmd.Flags().GetString("herokuApp")
	if err != nil {
		fmt.Println("error getting herokuApp flag", err)
		os.Exit(1)
	}

	deployerConfig, err := config.NewDeployerConfig(config.DeployerConfigFile, herokuApp)
	if err != nil {
		logger.Fatalln("error getting deployer config:", err)
	}

	if all {
		logger.Println("Uninstalling all apps")
		err := unInstallAll(deployerConfig.AppConfigs)
		if err != nil {
			logger.Fatalln("Error uninstalling all apps:", err)
		}
		logger.Println("Successfully uninstalled all apps")
		os.Exit(0)
	}

	if repoName == "" || manifestName == "" {
		logger.Fatalln("repoName and manifestName cannot be empty if not using the --all flag")
	}

	logger.Println(fmt.Sprintf("Uninstalling %s/%s", repoName, manifestName))
	err = unInstall(repoName, manifestName)
	if err != nil {
		logger.Fatalln(fmt.Sprintf("Error uninstalling %s/%s: %s", repoName, manifestName, err))
	}
	logger.Println(fmt.Sprintf("Successfully uninstalled %s/%s", repoName, manifestName))
}
