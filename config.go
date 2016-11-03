package RPRS

import (
	"fmt"
	"github.com/spf13/viper"
)

func loadDefaultSettings() {
	viper.SetDefault("LogPath", "logs/")
	viper.SetDefault("Port", ":1323")
}

func initializeConfig() error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("RPRS")
	viper.SetConfigName("config.json") // name of config file
	viper.AddConfigPath(".")           // optionally look for config in the working directory
	loadDefaultSettings()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	return nil
}
