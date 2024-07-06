package config

import (
	"github.com/spf13/viper"
	"log"
)

// struct for configuration
type conf struct {
	TORCONTROL_PASSWORD string
	TORSERVER_ADDRESS   string
	TORCONTROL_ADDRESS  string
}

var Cfg *conf

func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatal(err)
	}
	return
}
