package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Port           string        `mapstructure:"PORT"`
	GinMode        string        `mapstructure:"GIN_MODE"`
	AllowedOrigins string        `mapstructure:"ALLOWED_ORIGINS"`
	RequestTimeout time.Duration `mapstructure:"REQUEST_TIMEOUT"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	config := &Config{
		// Default values
		Port:           "8080",
		GinMode:        "release",
		AllowedOrigins: "http://localhost:3000",
		RequestTimeout: 30 * time.Second,
	}

	if err := viper.ReadInConfig(); err != nil {
		return config, nil // Return default config if no .env file
	}

	err := viper.Unmarshal(config)
	return config, err
}
