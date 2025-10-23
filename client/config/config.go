package config

import "github.com/komadiina/spelltext/shared/config"

type Config struct {
	NatsURL         string `yaml:"nats_url" env:"NATS_URL" env-default:"nats://127.0.0.1:4222"`
	MaxPollInterval int    `yaml:"max_wait_time" env:"MAX_WAIT_TIME" env-default:"5"`
	AudioEnabled    bool   `yaml:"audio_enabled" env:"AUDIO_ENABLED" env-default:"false"`
	ChatPort        int    `yaml:"chat_port" env:"CHAT_PORT" env-default:"50051"`
	StorePort       int    `yaml:"store_port" env:"STORE_PORT" env-default:"50052"`
	InventoryPort   int    `yaml:"inventory_port" env:"INVENTORY_PORT" env-default:"50053"`
	CharacterPort   int    `yaml:"character_port" env:"CHARACTER_PORT" env-default:"50054"`
	GambaPort       int    `yaml:"gamba_port" env:"GAMBA_PORT" env-default:"50055"`
	AuthPort        int    `yaml:"auth_port" env:"AUTH_PORT" env-default:"50056"`
	CombatPort      int    `yaml:"combat_port" env:"COMBAT_PORT" env-default:"50057"`
	BuildPort       int    `yaml:"build_port" env:"BUILD_PORT" env-default:"50058"`
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
