package config

import (
	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	AppEnv string
	Port   string
	App struct {
		Name string
		Env  string
		Port int
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}
}

var AppConfig Config

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic("Error loading config: " + err.Error())
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		panic("Error unmarshaling config: "+ err.Error())
	}
}
