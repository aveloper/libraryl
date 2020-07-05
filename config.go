package main

import (
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

//AppConfig - Server Config Parameters
type AppConfig struct {
	PORT       int
	MongoDBUri string

	Environment string
	Production  bool
}

var appConfig *AppConfig
var configOnce sync.Once

func getConfig() *AppConfig {
	configOnce.Do(func() {
		initConfig()
	})
	return appConfig
}

func initConfig() {
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Err in reading config: " + err.Error())
	}

	readConfig()
}

func readConfig() {
	env := viper.GetString("ENVIRONMENT")

	appConfig = &AppConfig{
		PORT:        viper.GetInt("PORT"),
		MongoDBUri:  viper.GetString("MONGO_DB_URI"),
		Environment: env,
		Production:  strings.EqualFold(env, "PRODUCTION"),
	}
	verifyConfig()
}

func verifyConfig() {
	if appConfig.PORT == 0 {
		panic("PORT is not set")
	}
	if appConfig.MongoDBUri == "" {
		panic("MONGO_DB_URI is not set")
	}
	if !appConfig.Production {
		log.Info("Server running in DEVELOPMENT mode")
	}
}
