package pgconfig

import (
	"github.com/komadiina/spelltext/shared/config"
)

type Config struct {
	Host string `yaml:"host" env:"PG_HOST" env-default:"postgres"`
	Port int    `yaml:"port" env:"PG_PORT" env-default:"5432"`
}

func LoadConfig() (*Config, error) {
	var cfg struct {
		Root Config `yaml:"postgres"`
	}

	err := config.LoadConfig(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg.Root, nil
}
