package util

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Environment    string `mapstructure:"ENVIRONMENT"`
	ArangoUser     string `mapstructure:"ARANGO_USER"`
	ArangoPassword string `mapstructure:"ARANGO_PASSWORD"`
	ArangoHost     string `mapstructure:"ARANGO_HOST"`
	ArangoDBName   string `mapstructure:"ARANGO_DB_NAME"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = viper.Unmarshal(&config)
	return
}
