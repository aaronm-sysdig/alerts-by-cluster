package configuration

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

type ConfigManager struct {
	config *Config
	log    *logrus.Logger
}

func NewConfigManager(log *logrus.Logger) *ConfigManager {
	return &ConfigManager{
		config: &Config{},
		log:    log,
	}
}

func (cm *ConfigManager) LoadConfig() error {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // or viper.SetConfigType("yml")
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.AutomaticEnv()          // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		cm.log.Printf("Config file (%s) not found, continuing without", viper.ConfigFileUsed())
	} else {
		cm.log.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Bind environment variables
	viper.BindEnv("secure_url")
	viper.BindEnv("secure_api_token")

	// Unmarshal the config into the Config struct
	err := viper.Unmarshal(cm.config)
	if err != nil {
		return err
	}

	return nil
}

func (cm *ConfigManager) ValidateConfig() error {
	if cm.config == nil {
		return errors.New("config is nil")
	}
	if cm.config.SecureURL == "" {
		return errors.New("missing SECURE_URL")
	}
	if cm.config.SecureAPIToken == "" {
		return errors.New("missing SECURE_API_TOKEN")
	}
	return nil
}

func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}
