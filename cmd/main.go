package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	db, err := repository.NewDB(cfg.DBFile)
	if err != nil {
		log.Fatalf("Error connect to repository: %+v", err)
	}
	defer db.Close()

	// init repository
	repo := repository.NewWalletRepo(db)

	// init usecases
	taskUC := usecases.NewWalletUsecase(repo)
	taskHandler := api.NewTaskHandler(taskUC)

	r := chi.NewRouter()
	r.Post("/api/task", taskHandler.CreateTask)
	r.Get("/api/tasks", taskHandler.GetTask)

	serverAddress := fmt.Sprintf("localhost:%s", cfg.Port)
	log.Infoln("Listening on " + serverAddress)
	if err = http.ListenAndServe(serverAddress, r); err != nil {
		log.Panicf("Start server error: %+v", err.Error())
	}
}
