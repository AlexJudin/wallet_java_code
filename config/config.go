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
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
	}

	logLevel, err := log.ParseLevel(os.Getenv("LOGLEVEL"))
	if err != nil {
		return nil, err
	}

	cfg.LogLevel = logLevel

	dbCfg := СonfigDB{
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	cfg.СonfigDB = &dbCfg

	return &cfg, nil
}

func (c *Сonfig) GetDataSourceName() string {
	str := fmt.Sprintf("host=db port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.СonfigDB.Port, c.СonfigDB.User, c.СonfigDB.Password, c.СonfigDB.DBName)

	return str
}
