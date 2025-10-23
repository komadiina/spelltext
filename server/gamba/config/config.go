package config

import "github.com/komadiina/spelltext/shared/config"

type Config struct {
	ServicePort          int    `yaml:"port" env:"PORT" env-default:"50055"`
	PgUser               string `yaml:"pgUser" env:"PG_USER" env-default:"postgres"`
	PgPass               string `yaml:"pgPass" env:"PG_PASS" env-default:"changeme"`
	PgHost               string `yaml:"pgHost" env:"PG_HOST" env-default:"spelltext-postgresql-ha-pgpool.spelltext.svc.cluster.local"`
	PgPort               int    `yaml:"pgPort" env:"PG_PORT" env-default:"5432"`
	PgDbName             string `yaml:"pgDbName" env:"PG_DB_NAME" env-default:"spelltext"`
	PgSSLMode            string `yaml:"pgSslMode" env:"PG_SSL_MODE" env-default:"disable"`
	InventoryServicePort int    `yaml:"inventoryserver.port" env:"INVENTORYSERVER_PORT" env-default:"50053"`
}

func LoadConfig() (*Config, error) {
	var cfg struct {
		Root     Config `yaml:"gambaserver"`
	}

	err := config.LoadConfig(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg.Root, nil
}
