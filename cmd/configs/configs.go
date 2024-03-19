package configs

import (
	"github.com/spf13/viper"
)

type conf struct {
	RedisHost          string `mapstructure:"DB_HOST"`
	RequestsByIp       int    `mapstructure:"REQUESTS_BY_IP"`
	RequestsByToken    int    `mapstructure:"REQUESTS_BY_TOKEN"`
	TimeBlockedByIp    int    `mapstructure:"TIME_BLOCKED_BY_IP"`
	TimeBlockedByToken int    `mapstructure:"TIME_BLOCKED_BY_TOKEN"`
	WebServerPort      string `mapstructure:"WEB_SERVER_PORT"`
	TokenAllowed       string `mapstructure:"TOKEN_ALLOWED"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg, nil
}
