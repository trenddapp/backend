package etherscan

import "go.uber.org/config"

type Config struct {
	APIKey string `yaml:"apiKey"`
}

func NewConfig(cfg *config.YAML) (*Config, error) {
	c := &Config{}
	if err := cfg.Get("client.etherscan").Populate(c); err != nil {
		return nil, err
	}

	return c, nil
}
