package kafka

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"

	"github.com/AlexJudin/wallet_java_code/internal/model"
	"github.com/AlexJudin/wallet_java_code/internal/usecases"
)

const (
	PaymentOperationTopic = "PaymentOperation"
)

type PaymentOperationTopicHandler struct {
	uc           usecases.Wallet
	partConsumer sarama.PartitionConsumer
}

func NewPaymentOperationTopicHandler(uc usecases.Wallet, partConsumer sarama.PartitionConsumer) *PaymentOperationTopicHandler {
	return &PaymentOperationTopicHandler{
		uc:           uc,
		partConsumer: partConsumer,
	}
}

func (c PaymentOperationTopicHandler) CreateOperation(id int, partConsumer sarama.PartitionConsumer) {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case msg, ok := <-partConsumer.Messages():
			if !ok {
				log.Println("Channel closed, exiting goroutine")
				return
			}

			var paymentOperation model.PaymentOperation
			err := json.Unmarshal(msg.Value, &paymentOperation)

			if err != nil {
				log.Printf("Error unmarshaling JSON: %+v", err)
				continue
			}
			err = c.uc.CreateOperation(&paymentOperation)
			if err != nil {
				log.Printf("Error creating operation: %+v", err)
				continue
			}
		case <-stop:
			log.Println("Stop signal received. Shutting down kafka server...")
			partConsumer.Close()
			log.Info("The kafka server has been stopped successfully")
		}
	}
}
