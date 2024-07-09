package configuration

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
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

	// Define command-line flags
	pflag.String("secure_url", "", "Secure URL for the application")
	pflag.String("secure_api_token", "", "Secure API token for the application")
	pflag.Parse()

	// Bind command-line flags to Viper
	viper.BindPFlag("secure_url", pflag.Lookup("secure_url"))
	viper.BindPFlag("secure_api_token", pflag.Lookup("secure_api_token"))

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
