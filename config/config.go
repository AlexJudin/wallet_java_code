package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Сonfig struct {
	Host     string
	Port     string
	LogLevel log.Level
	СonfigDB *СonfigDB
}

type СonfigDB struct {
	Port     string
	Host     string
	User     string
	Password string
	DBName   string
	Sslmode  string
}

func New() (*Сonfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	cfg := Сonfig{
		Port: os.Getenv("PORT"),
	}

	logLevel, err := log.ParseLevel(os.Getenv("LOGLEVEL"))
	if err != nil {
		return nil, err
	}

	cfg.LogLevel = logLevel

	return &cfg, nil
}

func (c *Сonfig) GetDataSourceName() string {
	str := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.СonfigDB.Host, c.СonfigDB.Port, c.СonfigDB.User, c.СonfigDB.Password, c.СonfigDB.DBName, c.СonfigDB.Sslmode)

	return str
}
