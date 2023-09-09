package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environtment        string        `mapstructure:"ENVIRONTMENT"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	AccessTokenKey      string        `mapstructure:"ACCESS_TOKEN_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	MidtransServerKey   string        `mapstructure:"MIDTRANS_SERVER_KEY"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	// viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	// opt := viper.DecodeHook(mapstructure.StringToTimeDurationHookFunc())

	err = viper.Unmarshal(&config)
	return
}
