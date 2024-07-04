package configuration

type Config struct {
	SecureURL      string `mapstructure:"secure_url"`
	SecureAPIToken string `mapstructure:"secure_api_token"`
}
