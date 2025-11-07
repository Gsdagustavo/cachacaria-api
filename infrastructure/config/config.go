package config

import (
	"cachacariaapi/infrastructure/modules"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
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

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	if path == "" {
		path = "./build/config/dev.toml"
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

func (s Server) RegisterModules(router *mux.Router, modules ...modules.Module) {
	for _, module := range modules {
		module.RegisterRoutes(router)
	}
}

func (s Server) Run(cfg Config) error {
	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting HTTP server on %s", address)

	//s.Router.Use(LoggingMiddleware)
	s.Router.Use(CORSMiddleware)

	s.Router.PathPrefix("/images/").Handler(http.StripPrefix("/images/",
		http.FileServer(http.Dir("/app/images")),
	))

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
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
	)
}

func Connect(cfg *Config) error {
	log.Printf("Connecting to database %s", cfg.Database.Name)
	log.Printf("Driver: %s", cfg.Database.Driver)
	var conn *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		conn, err = sql.Open(cfg.Database.Driver, cfg.Database.GetDSN())
		if err != nil {
			log.Printf("Error connecting to database: %s", err)
			log.Printf("Retrying in 1 second...")
			time.Sleep(1 * time.Second)
			continue
		}
	}

	cfg.Database.Conn = conn
	log.Println("Database connection established successfully")
	return nil
}

//func LoggingMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		slog.Info(
//			"Incoming request",
//			"method",
//			r.Method,
//			"url",
//			r.URL.Path,
//			"remote_addr",
//			r.RemoteAddr,
//			"user_agent",
//			r.UserAgent(),
//			"host",
//			r.Host,
//			"cookies",
//			r.Cookies(),
//			"body",
//			r.Body,
//		)
//		next.ServeHTTP(w, r)
//	})
//}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		if r.Method == http.MethodOptions {
			log.Printf("[CORS Middleware] allow options | no content")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
