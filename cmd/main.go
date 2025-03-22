package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"

	"github.com/AlexJudin/wallet_java_code/api"
	"github.com/AlexJudin/wallet_java_code/config"
	"github.com/AlexJudin/wallet_java_code/repository"
	"github.com/AlexJudin/wallet_java_code/usecases"
)

// @title Пользовательская документация API
// @description Тестовое задание
// @termsOfService spdante@mail.ru
// @contact.name Alexey Yudin
// @contact.email spdante@mail.ru
// @version 1.0.0
// @host localhost:7540
// @BasePath /
func main() {
	// init config
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(cfg.LogLevel)

	connStr := cfg.GetDataSourceName()
	db, err := repository.ConnectDB(connStr)
	if err != nil {
		log.Fatal(err)
	}

	// init repository
	repo := repository.NewWalletRepo(db)

	// init usecases
	walletUC := usecases.NewWalletUsecase(repo)
	walletHandler := api.NewWalletHandler(walletUC)

	r := chi.NewRouter()
	r.Post("/api/v1/wallet", walletHandler.CreateOperation)
	r.Get("/api/v1/wallets/", walletHandler.GetWalletBalanceByUUID)

	log.Info("Start http server")

	serverAddress := fmt.Sprintf(":%s", cfg.Port)
	log.Infoln("Listening on " + serverAddress)
	if err = http.ListenAndServe(serverAddress, r); err != nil {
		log.Panicf("Start server error: %+v", err.Error())
	}
}
