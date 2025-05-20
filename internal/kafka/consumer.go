package kafka

import (
	"github.com/IBM/sarama"

	"github.com/AlexJudin/wallet_java_code/config"
)

func NewConsumer(cfg *config.Сonfig) (sarama.Consumer, error) {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
