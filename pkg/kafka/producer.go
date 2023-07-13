package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type kafkaProducer struct {
	producer sarama.SyncProducer
}

type KafkaProducer interface {
	PublishMessage(topic string, message []byte) error
}

func ConnectProducer(brokersUrl []string) (KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	// NewSyncProducer creates a new SyncProducer using the given broker addresses and configuration.
	conn, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	producer := &kafkaProducer{
		producer: conn,
	}
	return producer, nil
}

func (p *kafkaProducer) PublishMessage(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := p.producer.SendMessage(msg)

	if err != nil {
		return err
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)

	return nil
}
