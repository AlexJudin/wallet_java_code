package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Сonfig struct {
	Port     string
	LogLevel log.Level
	СonfigDB *СonfigDB
}

type СonfigDB struct {
	Port     string `env:"DBPORT" envDefault:"5433"`
	Host     string `env:"DBHOST,required"`
	User     string `env:"DBUSER" envDefault:"http"`
	Password string `env:"DBPASSWORD,required"`
	DBName   string `env:"DBNAME" envDefault:"http"`
	Sslmode  string `env:"SSLMODE" envDefault:"disable"`
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

func (c *Сonfig) GetDataSourceName() string {
	str := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.СonfigDB.Host, c.СonfigDB.Port, c.СonfigDB.User, c.СonfigDB.Password, c.СonfigDB.DBName, c.СonfigDB.Sslmode)

	return str
}
