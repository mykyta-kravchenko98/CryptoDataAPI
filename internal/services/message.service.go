package services

import (
	"encoding/json"
	"errors"

	"github.com/mykyta-kravchenko98/CryptoDataAPI/internal/configs"
	"github.com/mykyta-kravchenko98/CryptoDataAPI/internal/models"
	"github.com/mykyta-kravchenko98/CryptoDataAPI/pkg/kafka"
	"github.com/rs/zerolog/log"
)

type messageService struct {
	kafkaProducer kafka.KafkaProducer
	configs       configs.KafkaConfig
}

type MessageService interface {
	SendCryptoDataMessage(coins []models.Coin) error
}

func NewMessageService(c configs.KafkaConfig) MessageService {
	producer, err := kafka.ConnectProducer([]string{c.URL})

	if err != nil {
		log.Fatal().Err(err)
	}
	instance := &messageService{
		kafkaProducer: producer,
		configs:       c,
	}

	return instance
}

func (ms *messageService) SendCryptoDataMessage(coins []models.Coin) error {
	topicName := ms.configs.CryptodataTopic
	if topicName == "" {
		return errors.New("CryptodataTopic config is empty")
	}

	data, err := json.Marshal(coins)
	if err != nil {
		return err
	}

	err = ms.kafkaProducer.PublishMessage(topicName, data)

	return err
}
