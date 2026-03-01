package configs

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	Port     string
	DBConfig DBConfig
}

type DBConfig struct {
	DBUrl string
}

func LoadAppConfig() AppConfig {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic("Failed to read config file: " + err.Error())
	}

	appConfig := AppConfig{
		Port: viper.GetString("PORT"),
		DBConfig: DBConfig{
			DBUrl: viper.GetString("DB_URL"),
		},
	}

	return appConfig
}
