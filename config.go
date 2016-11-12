package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func loadDefaultSettings() {
	viper.SetDefault("LogPath", "logs/")
	viper.SetDefault("Port", ":1323")
	viper.SetDefault("UploadRpmPath", "uploads/")
	viper.SetDefault("Secret", "secret")
}

func initializeConfig() error {
	//viper.AutomaticEnv()
	viper.SetConfigType("json")
	viper.SetEnvPrefix("RPRS")
	viper.SetConfigName("config") // name of config file
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	loadDefaultSettings()
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	return nil
}
