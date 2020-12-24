package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var config *viper.Viper

func Init() {
	config = viper.New()
	config.SetConfigType("yml")
	config.SetConfigName("config")
	config.AddConfigPath("/app")
	config.AddConfigPath(".")
	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func GetConfig() *viper.Viper {
	return config
}
