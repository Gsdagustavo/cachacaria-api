package config

import (
	"cachacariaapi/infrastructure/modules"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
)

type Config struct {
	Database     Database `toml:"Database"`
	Server       Server   `toml:"Server"`
	SymmetricKey string   `toml:"symmetric_key"`
}

func LoadConfig() (*Config, error) {
	var path string
	flag.StringVar(&path, "config", "", "Path to config TOML file")
	flag.Parse()

	// fallback to env var if not provided
	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	if path == "" {
		path = "./build/config/dev.toml" // default fallback
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

type Server struct {
	Port    int    `toml:"port"`
	Host    string `toml:"host"`
	BaseURL string `toml:"base_url"`
	Router  *mux.Router
	Server  *http.Server
}

func (s Server) RegisterModules(modules ...modules.Module) {
	for _, module := range modules {
		module.RegisterRoutes(s.Router)
	}
}

func (s Server) Run(cfg Config) error {
	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting HTTP server on %s", address)

	s.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))
	return http.ListenAndServe(address, s.Router)
}

type Database struct {
	Driver   string `toml:"driver"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Name     string `toml:"name"`
	Conn     *sql.DB
}

func (d Database) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", d.User, d.Password, d.Host, d.Port, d.Name)
}

func Connect(cfg *Config) error {
	log.Printf("Connecting to database %s", cfg.Database.Name)
	log.Printf("Driver: %s", cfg.Database.Driver)

	conn, err := sql.Open(cfg.Database.Driver, cfg.Database.GetDSN())
	if err != nil {
		return fmt.Errorf("error opening connection: %s", err)
	}

	err = conn.Ping()
	if err != nil {
		return fmt.Errorf("error pinging database: %s", err)
	}

	cfg.Database.Conn = conn
	log.Println("Database connection established successfully")
	return nil
}
