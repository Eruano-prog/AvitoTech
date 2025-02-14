package config

import (
	_ "embed"
	"gopkg.in/yaml.v3"
)

// TODO: replace me with .env
//
//go:embed config.yaml
var cfg []byte

type databaseConfig struct {
	Postgres struct {
		Address  string `yaml:"address"`
		DBName   string `yaml:"db_name"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"postgres"`
}

type serverConfig struct {
	RESTAddr string `yaml:"restAddr"`
}

type Config struct {
	JwtSecret string         `yaml:"jwtSecret"`
	Database  databaseConfig `yaml:"database"`
	Server    serverConfig   `yaml:"server"`
}

var Configuration = Config{}

// LoadConfiguration expected to refresh cfgs from file but now it just initializes it
// TODO: Seems like it useless because the file is precompiled with config. Need to separate static part and mutual and rewrite it without go:embed to provide possibility to refresh config from file
func LoadConfiguration() error {
	return yaml.Unmarshal(cfg, &Configuration)
}
