package config

import "github.com/spf13/viper"

type EnvConfig struct {
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort int    `mapstructure:"REDIS_PORT"`
	RedisPass string `mapstructure:"REDIS_PWD"`

	RateLimitBy        string `mapstructure:"RATE_LIMIT_BY"`
	IpLimitRate        int    `mapstructure:"IP_LIMIT_RATE"`
	TokenLimitRate     int    `mapstructure:"TOKEN_LIMIT_RATE"`
	IpWindowTime       int    `mapstructure:"IP_WINDOW_TIME"`
	ApiTokenWindowTime int    `mapstructure:"API_TOKEN_WINDOW_TIME"`

	AppPort string `mapstructure:"APP_PORT"`
}

func SetupEnvConfig() (*EnvConfig, error) {
	viper.SetConfigName("env-rl")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var envConf EnvConfig
	if err := viper.Unmarshal(&envConf); err != nil {
		return nil, err
	}

	return &envConf, nil
}
