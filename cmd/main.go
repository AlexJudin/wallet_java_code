package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
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
	r.Use(httprate.LimitByIP(1000, time.Second))
	r.Post("/api/v1/wallet", walletHandler.CreateOperation)
	r.Get("/api/v1/wallets/", walletHandler.GetWalletBalanceByUUID)

	log.Info("Start http server")

	serverAddress := fmt.Sprintf(":%s", cfg.Port)
	serverErr := make(chan error)

	httpServer := &http.Server{
		Addr:    serverAddress,
		Handler: r,
	}

	go func() {
		log.Infof("Listening on %s", serverAddress)
		if err = httpServer.ListenAndServe(); err != nil {
			serverErr <- err
		}
		close(serverErr)
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		log.Info("Stop signal received. Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = httpServer.Shutdown(ctx); err != nil {
			log.Errorf("error terminating server: %+v", err)
		}
		log.Info("The server has been stopped successfully")
	case err = <-serverErr:
		log.Errorf("Server error: %+v", err)
	}
}
