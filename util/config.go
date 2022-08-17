package util

import (
	"time"

	"github.com/spf13/viper"
)

//Config contiene tutti i valori delle configurazioni che vengono usate nella app
//e che leggiamo dal file di config
type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_URL"`
	ServerHost          string        `mapstructure:"SERVER_HOST"`
	ServerPort          string        `mapstructure:"SERVER_PORT"`
	GinRelease          string        `mapstructure:"GIN_MODE"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SIMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(configPath string) (config Config, err error) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
