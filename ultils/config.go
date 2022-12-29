package ultils

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddr  string `mapstructure:"SERVER_ADDRESS"`
	Environment string `mapstructure:"ENVIRONMENT"`
}

func LoadConfig() (Config, error) {
	var config Config
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}
	viper.Unmarshal(&config)

	return config, nil
}
