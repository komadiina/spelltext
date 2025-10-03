package config

import "github.com/komadiina/spelltext/shared/config"

type Config struct {
	Port            int    `yaml:"port" env:"PORT" env-default:"50051"`
	NatsURL         string `yaml:"natsUrl" env:"NATS_URL" env-default:"nats://spelltext-nats:4222"`
	MaxAsyncPublish int    `yaml:"maxAsyncPublish" env:"MAX_ASYNC_PUBLISH" env-default:"512"`
}

func LoadConfig() (*Config, error) {
	var cfg struct {
		Root Config `yaml:"chatserver"`
	}

	err := config.LoadConfig(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg.Root, nil
}
