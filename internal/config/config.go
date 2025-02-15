// Package config provides config struct which are should be loaded from .env
package config

type Config struct {
	JwtSecret string `env:"JWT_SECRET" env-required:"true"`
	ItemsPath string `env:"ITEMS_PATH" env-required:"true"`
	Database  databaseConfig
	Server    serverConfig
}

type databaseConfig struct {
	Address  string `env:"DATABASE_ADDRESS" env-required:"true"`
	DBName   string `env:"DATABASE_DB_NAME" env-required:"true"`
	Username string `env:"DATABASE_USERNAME" env-required:"true"`
	Password string `env:"DATABASE_PASSWORD" env-required:"true"`
}

type serverConfig struct {
	RESTAddr string `env:"SERVER_REST_ADDR" env-required:"true"`
}

var Configuration Config
