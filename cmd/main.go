package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/AlexJudin/wallet_java_code/api"
	"github.com/AlexJudin/wallet_java_code/config"
	"github.com/AlexJudin/wallet_java_code/repository"
	"github.com/AlexJudin/wallet_java_code/usecases"
)

// @title Пользовательская документация API
// @description Итоговая работа по курсу "Go-разработчик с нуля" (Яндекс Практикум)
// @termsOfService spdante@mail.ru
// @contact.name Alexey Yudin
// @contact.email spdante@mail.ru
// @version 1.0.0
// @host localhost:7540
// @BasePath /
func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(cfg.LogLevel)

	connStr := cfg.GetDataSourceName()
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		cfg.СonfigDB.DBName, driver)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}

	// init repository
	repo := repository.NewWalletRepo(db)

	// init usecases
	walletUC := usecases.NewWalletUsecase(repo)
	walletHandler := api.NewWalletHandler(walletUC)

	r := chi.NewRouter()
	r.Post("/api/v1/wallet", walletHandler.CreateOperation)
	r.Get("/api/v1/wallets/{WALLET_UUID:string}", walletHandler.GetWalletByUUID)

	serverAddress := fmt.Sprintf("localhost:%s", cfg.Port)
	log.Infoln("Listening on " + serverAddress)
	if err = http.ListenAndServe(serverAddress, r); err != nil {
		log.Panicf("Start server error: %+v", err.Error())
	}
}
