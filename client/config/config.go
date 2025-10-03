package config

import "github.com/komadiina/spelltext/shared/config"

type Config struct {
	NatsURL         string `yaml:"nats_url" env:"NATS_URL" env-default:"nats://127.0.0.1:4222"`
	MaxPollInterval int    `yaml:"max_wait_time" env:"MAX_WAIT_TIME" env-default:"5"`
}

func LoadConfig() (*Config, error) {
	var cfg struct {
		Root Config `yaml:"client"`
	}

	err := config.LoadConfig(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg.Root, nil
}
