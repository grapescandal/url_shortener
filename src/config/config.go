package config

import "github.com/spf13/viper"

type AppServerConfig struct {
	AppName  string `mapstructure:"APP_NAME"`
	HTTPHost string `mapstructure:"HTTP_HOST"`
}

func LoadConfig() (config *AppServerConfig, err error) {
	viper.SetConfigFile("devops/.env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
