package configs

import (
	"github.com/spf13/viper"
)

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
