package configuration

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

func LoadConfig(log *logrus.Logger) (*Config, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // or viper.SetConfigType("yml")
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.AutomaticEnv()          // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values if needed
	viper.SetDefault("secure_url", "https://default.url")
	viper.SetDefault("secure_api_token", "default-token")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %s", err)
	} else {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Bind environment variables
	viper.BindEnv("secure_url")
	viper.BindEnv("secure_api_token")

	// Unmarshal the config into the Config struct
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
