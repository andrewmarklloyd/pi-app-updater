package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/andrewmarklloyd/pi-app-deployer/internal/pkg/config"
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
	herokuAPIKey := os.Getenv("HEROKU_API_KEY")
	if herokuAPIKey == "" {
		logger.Fatalln("HEROKU_API_TOKEN environment variable is required")
	}

	_, err := newAgent(herokuAPIKey)
	if err != nil {
		logger.Fatalln(fmt.Errorf("error creating agent: %s", err))
	}

	appConfigs, err := config.GetAppConfigs(config.AppConfigsFile)
	if err != nil {
		logger.Fatalln("error getting app configs:", err)
	}

	fmt.Println(appConfigs)
	// if !appConfigs.ConfigExists(cfg) {
	// 	logger.Fatalln("App does not exist in app configs file:", config.AppConfigsFile)
	// }

	// err = agent.MqttClient.Connect()
	// if err != nil {
	// 	logger.Fatalln("connecting to mqtt: ", err)
	// }

	// updateCondition := config.UpdateCondition{
	// 	RepoName:     cfg.RepoName,
	// 	ManifestName: cfg.ManifestName,
	// }

	// agent.MqttClient.Subscribe(config.RepoPushTopic, func(message string) {
	// 	var artifact config.Artifact
	// 	err := json.Unmarshal([]byte(message), &artifact)
	// 	if err != nil {
	// 		logger.Println(fmt.Sprintf("unmarshalling payload from topic %s: %s", config.RepoPushTopic, err))
	// 		return
	// 	}
	// 	if artifact.RepoName == cfg.RepoName && artifact.ManifestName == cfg.ManifestName {
	// 		logger.Println(fmt.Sprintf("updating repo %s with manifest name %s", cfg.RepoName, cfg.ManifestName))
	// 		updateCondition.Status = config.StatusInProgress
	// 		err = agent.publishUpdateCondition(updateCondition)
	// 		if err != nil {
	// 			// log but don't block update from proceeding
	// 			logger.Println(err)
	// 		}
	// 		err := agent.handleRepoUpdate(artifact)
	// 		if err != nil {
	// 			logger.Println(err)
	// 			updateCondition.Status = config.StatusErr
	// 			err = agent.publishUpdateCondition(updateCondition)
	// 			if err != nil {
	// 				logger.Println(err)
	// 			}
	// 			return
	// 		}
	// 		// TODO: should check systemctl status before sending success?
	// 		updateCondition.Status = config.StatusSuccess
	// 		err = agent.publishUpdateCondition(updateCondition)
	// 		if err != nil {
	// 			logger.Println(err)
	// 		}
	// 	}
	// })

	// agent.MqttClient.Subscribe(config.ServiceActionTopic, func(message string) {
	// 	var payload config.ServiceActionPayload
	// 	err := json.Unmarshal([]byte(message), &payload)
	// 	if err != nil {
	// 		logger.Println(fmt.Sprintf("unmarshalling payload from topic %s: %s", config.ServiceActionTopic, err))
	// 		return
	// 	}
	// 	if payload.RepoName == cfg.RepoName && payload.ManifestName == cfg.ManifestName {
	// 		logger.Println(fmt.Sprintf("Running service action %s on %s/%s", payload.Action, payload.RepoName, payload.ManifestName))
	// 		var err error
	// 		switch payload.Action {
	// 		case config.ServiceActionStart:
	// 			err = file.StartSystemdUnit(payload.ManifestName)
	// 			break
	// 		case config.ServiceActionStop:
	// 			err = file.StopSystemdUnit(payload.ManifestName)
	// 			break
	// 		case config.ServiceActionRestart:
	// 			err = file.RestartSystemdUnit(payload.ManifestName)
	// 			break
	// 		default:
	// 			err = fmt.Errorf("Action %s is not valid", payload.Action)
	// 			break
	// 		}
	// 		if err != nil {
	// 			logger.Println(err)
	// 		}

	// 	}
	// })

	// config.AppConfigs
	// if *logForwarding {
	// 	logger.Println(fmt.Sprintf("Log forwarding is enabled for %s", cfg.ManifestName))
	// 	agent.startLogForwarder(cfg.ManifestName, func(log string) {
	// 		l := config.Log{
	// 			Message: log,
	// 			Config:  cfg,
	// 		}
	// 		json, err := json.Marshal(l)
	// 		if err != nil {
	// 			logger.Println(fmt.Sprintf("marshalling log forwarder message: %s", err))
	// 			return
	// 		}
	// 		err = agent.MqttClient.Publish(config.LogForwarderTopic, string(json))
	// 		if err != nil {
	// 			logger.Println(fmt.Sprintf("error publishing log forwarding message: %s", err))
	// 		}
	// 	})
	// }

	go forever()
	select {} // block forever

}

func forever() {
	for {
		time.Sleep(5 * time.Minute)
	}
}
