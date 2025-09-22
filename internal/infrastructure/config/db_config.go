package config

import (
	"os"
	"strconv"
)

type DBConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func NewDBConfig() *DBConfig {
	var (
		driver   = os.Getenv("DB_DRIVER")
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		passwd   = os.Getenv("DB_PASSWORD")
		database = os.Getenv("DB_NAME")
	)

	intPort, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}

	return &DBConfig{
		Driver:   driver,
		Host:     host,
		Port:     intPort,
		User:     user,
		Password: passwd,
		Database: database,
	}
}

func (c *DBConfig) FormatDSN() string {
	cred := ""
	if c.User != "" {
		cred = c.User
		if c.Password != "" {
			cred += ":" + c.Password
		}
	}
	addr := c.Host
	if c.Port > 0 {
		addr += ":" + strconv.Itoa(c.Port)
	}
	dsn := ""
	if cred != "" {
		dsn = cred + "@"
	}
	dsn += "tcp(" + addr + ")/"
	dsn += c.Database

	return dsn
}
