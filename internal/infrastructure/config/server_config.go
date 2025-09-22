package config

import (
	"os"
	"strings"
)

type ServerConfig struct {
	Port      string
	Address   string
	BaseURL   string
	JWTSecret string
	DBConfig  *DBConfig
}

func NewServerConfig(db *DBConfig) *ServerConfig {
	if db == nil {
		db = NewDBConfig()
	}

	port := os.Getenv("SERVER_PORT")
	addr := os.Getenv("SERVER_ADDRESS")
	baseURL := os.Getenv("BASE_URL")
	jwt := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if jwt == "" {
		panic("JWT_SECRET is required")
	}

	return &ServerConfig{
		Port:      port,
		Address:   addr,
		BaseURL:   baseURL,
		JWTSecret: jwt,
		DBConfig:  db,
	}
}
