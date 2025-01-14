package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Configuration for app
type Configuration struct {
	Domain              string `required:"true"`
	MaxAddr             int    `split_words:"true" required:"true"`
	RequestOverTime     int    `split_words:"true" required:"true"`
	CookieHashKey       string `split_words:"true" required:"true"`
	CookieBlockKey      string `split_words:"true" required:"true"`
	CoinMarketCapApiKey string `split_words:"true" required:"true"`
	SmtpHost            string `split_words:"true" required:"true"`
	SmtpPort            string `split_words:"true" required:"true"`
	SmtpUser            string `split_words:"true" required:"true"`
	SmtpPassword        string `split_words:"true" required:"true"`
}

func New() *Configuration {
	return &Configuration{}
}

// GetEnv configuration init
func (cnf *Configuration) GetEnv() error {
	if err := envconfig.Process("", cnf); err != nil {
		return err
	}
	return nil
}
