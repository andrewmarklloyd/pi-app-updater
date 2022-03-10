package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andrewmarklloyd/pi-app-deployer/internal/pkg/config"
	"github.com/andrewmarklloyd/pi-app-deployer/internal/pkg/mqtt"
	"github.com/andrewmarklloyd/pi-app-deployer/internal/pkg/redis"
	gmux "github.com/gorilla/mux"
)

var logger = log.New(os.Stdout, "[pi-app-deployer-Server] ", log.LstdFlags)
var forwarderLogger = log.New(os.Stdout, "[pi-app-deployer-Forwarder] ", log.LstdFlags)

var messageClient mqtt.MqttClient

func main() {
	srvAddr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))

	// TODO: reusing another app's mqtt instance to save cost. Once viable MVP finished I can provision a dedicated instance
	user := os.Getenv("CLOUDMQTT_USER")
	pw := os.Getenv("CLOUDMQTT_PASSWORD")
	url := os.Getenv("CLOUDMQTT_URL")
	mqttAddr := fmt.Sprintf("mqtt://%s:%s@%s", user, pw, url)

	messageClient = mqtt.NewMQTTClient(mqttAddr, *logger)
	err := messageClient.Connect()
	if err != nil {
		logger.Fatalln("connecting to mqtt: ", err)
	}

	redisClient, err := redis.NewRedisClient(os.Getenv("REDIS_TLS_URL"))
	if err != nil {
		logger.Fatalln("creating redis client:", err)
	}

	messageClient.Subscribe(config.LogForwarderTopic, func(message string) {
		var log config.Log
		err := json.Unmarshal([]byte(message), &log)
		if err != nil {
			logger.Println(fmt.Sprintf("unmarshalling log forwarder message: %s", err))
		}

		forwarderLogger.Println(fmt.Sprintf("<%s_%s>: %s", log.Config.RepoName, log.Config.ManifestName, log.Message))
	})

	messageClient.Subscribe(config.RepoPushStatusTopic, func(message string) {
		var c config.UpdateCondition
		err := json.Unmarshal([]byte(message), &c)
		if err != nil {
			logger.Println(fmt.Sprintf("unmarshalling update condition message: %s", err))
			return
		}
		logger.Println(fmt.Sprintf("<%s/%s> update status: %s", c.RepoName, c.ManifestName, c.Status))

		key := fmt.Sprintf("%s/%s", c.RepoName, c.ManifestName)
		value := c.Status
		err = redisClient.WriteCondition(context.Background(), key, value)
		if err != nil {
			logger.Println(fmt.Sprintf("writing condition message to redis: %s", err))
			return
		}
	})

	router := gmux.NewRouter().StrictSlash(true)
	router.Handle("/push", requireLogin(http.HandlerFunc(handleRepoPush))).Methods("POST")

	srv := &http.Server{
		Handler: router,
		Addr:    srvAddr,
	}

	logger.Println("server started on ", srvAddr)
	logger.Fatal(srv.ListenAndServe())
}

func requireLogin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if !isAuthenticated(req) {
			logger.Println(fmt.Sprintf("Unauthenticated request, host: %s, headers: %s", req.Host, req.Header))
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

func isAuthenticated(req *http.Request) bool {
	allowedApiKey := os.Getenv("PI_APP_DEPLOYER_API_KEY")
	apiKey := req.Header.Get("api-key")
	if apiKey == "" {
		return false
	}
	if apiKey != allowedApiKey {
		return false
	}
	return true
}
