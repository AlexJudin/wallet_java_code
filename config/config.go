package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const (
	accessTokenTTLDefault  = 5
	refreshTokenTTLDefault = 60
)

type Сonfig struct {
	Host     string
	Port     string
	LogLevel log.Level
	*СonfigDB
	*ConfigAuth
	*ConfigRedis
}

type СonfigDB struct {
	Port     string
	User     string
	Password string
	DBName   string
	Sslmode  string
}

type ConfigAuth struct {
	PasswordSalt    string
	TokenSalt       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type ConfigRedis struct {
	Host     string
	Port     string
	Password string
}

func New() (*Сonfig, error) {
	err := godotenv.Load("config/config.env")
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

	accessTokenTTL, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil {
		accessTokenTTL = accessTokenTTLDefault
	}

	refreshTokenTTL, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil {
		refreshTokenTTL = refreshTokenTTLDefault
	}

	authCfg := ConfigAuth{
		PasswordSalt:    os.Getenv("PASSWORD_SALT"),
		TokenSalt:       os.Getenv("TOKEN_SALT"),
		AccessTokenTTL:  accessTokenTTL * time.Minute,
		RefreshTokenTTL: refreshTokenTTL * time.Minute,
	}
	cfg.ConfigAuth = &authCfg

	redisCfg := ConfigRedis{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}
	cfg.ConfigRedis = &redisCfg

	return &cfg, nil
}

func (c *Сonfig) GetDataSourceName() string {
	str := fmt.Sprintf("host=db port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.СonfigDB.Port, c.СonfigDB.User, c.СonfigDB.Password, c.СonfigDB.DBName)

	return str
}
