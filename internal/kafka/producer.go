package kafka

import (
	"github.com/IBM/sarama"

	"github.com/AlexJudin/wallet_java_code/config"
)

func NewProducer(cfg *config.Сonfig) (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
