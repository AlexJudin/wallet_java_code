package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexJudin/wallet_java_code/internal/usecases"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"

	"github.com/AlexJudin/wallet_java_code/config"
	"github.com/AlexJudin/wallet_java_code/internal/api/controller"
	"github.com/AlexJudin/wallet_java_code/internal/cache"
	"github.com/AlexJudin/wallet_java_code/internal/kafka"
	"github.com/AlexJudin/wallet_java_code/internal/model"
	"github.com/AlexJudin/wallet_java_code/internal/repository"
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

	redisClient, err := cache.ConnectToRedis(cfg)
	if err != nil {
		log.Error("error connecting to redis")
	}

	consumer, err := kafka.NewConsumer(cfg)
	if err != nil {
		log.Fatalf("failed to create consumer: %+v", err)
	}
	defer consumer.Close()

	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		log.Fatalf("failed to create producer: %+v", err)
	}
	defer producer.Close()

	r := chi.NewRouter()
	controller.AddRoutes(cfg, db, redisClient, r)

	startKafka(consumer, producer, db, redisClient)

	startHTTPServer(cfg, r)
}

func startHTTPServer(cfg *config.Сonfig, r *chi.Mux) {
	var err error

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

func startKafka(consumer sarama.Consumer, producer sarama.SyncProducer, db *gorm.DB, redisClient *redis.Client) {
	partConsumer, err := consumer.ConsumePartition(kafka.PaymentOperationTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %+v", err)
	}

	publish(producer)

	repoWallet := repository.NewWalletRepo(db)
	cache := cache.NewCacheClientRepo(redisClient)

	uc := usecases.NewWalletUsecase(repoWallet, cache)
	walletKafka := kafka.NewPaymentOperationTopicHandler(uc, partConsumer)

	for w := 1; w <= 3; w++ {
		go walletKafka.CreateOperation(w, partConsumer)
	}
}

func publish(producer sarama.SyncProducer) {
	for i := 1; i <= 15; i++ {
		msg := model.PaymentOperation{
			WalletId:      strconv.Itoa(i),
			OperationType: model.Deposit,
			Amount:        3000 + int64(i),
		}

		res, _ := json.Marshal(msg)
		resp := &sarama.ProducerMessage{
			Topic: kafka.PaymentOperationTopic,
			Key:   sarama.StringEncoder(strconv.Itoa(i)),
			Value: sarama.StringEncoder(res),
		}

		_, _, err := producer.SendMessage(resp)
		if err != nil {
			log.Printf("Failed to produce message: %+v", err)
		}
	}
}
