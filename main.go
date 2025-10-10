package main

import (
	"cachacariaapi/infrastructure"
	"flag"
	"log/slog"
	"os"
)

func main() {
	setupLogger()22222
	cfgFilePath := getCFGFileFlag()

	infrastructure.Init(cfgFilePath)
}

func getCFGFileFlag() string {
	filepath := flag.String("config", "dev.toml", "Path to config file")
	flag.Parse()

	if filepath == nil || *filepath == "" {
		fallback := "dev.toml"
		return fallback
	}

	return *filepath
}

func setupLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
}
