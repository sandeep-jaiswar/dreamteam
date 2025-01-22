package config

import (
	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
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
	Jwt struct {
		SecretKey string
	}
	GroupCache struct {
		Name      string
		CacheSize int64
	}
	Profiling struct {
		Port int
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


func ProvideConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic("Error loading config: " + err.Error())
	}

	var appConfig Config
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		panic("Error unmarshaling config: "+ err.Error())
	}

	return &appConfig
}
