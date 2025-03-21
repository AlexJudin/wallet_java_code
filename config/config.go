package config

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Сonfig struct {
	Port     string
	LogLevel log.Level
}

func New() (*Сonfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	cfg := Сonfig{
		Port: os.Getenv("TODO_PORT"),
	}

	logLevel, err := log.ParseLevel(os.Getenv("TODO_LOGLEVEL"))
	if err != nil {
		return nil, err
	}

	cfg.LogLevel = logLevel

	return &cfg, nil
}
