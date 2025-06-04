package config

import "github.com/spf13/viper"

type Config struct {
	Port              string `mapstructure:"PORT"`
	AuthServiceClient string `mapstructure:"AUTH_SERVICE_CLIENT"`
}

func LoadConfig() (Config, error) {
	var config Config
	viper.AddConfigPath("./")
	viper.SetConfigFile("./pkg/config/envs/dev.env")
	//viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}
	err := viper.Unmarshal(&config)
	return config, err
}
