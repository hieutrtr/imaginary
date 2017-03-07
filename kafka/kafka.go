package kafka

import (
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

var producer sarama.SyncProducer

// Event interface
type Event interface {
	payloadBuild() (*sarama.ProducerMessage, error)
}

// UploadEvent is uploading image event structure
type UploadEvent struct {
	Topic string
	Oid   string
}

func (e *UploadEvent) payloadBuild() (*sarama.ProducerMessage, error) {
	if e.Topic == "" {
		return nil, ErrMessage
	}
	mess := &sarama.ProducerMessage{
		Topic: e.Topic,
		Key:   sarama.StringEncoder(e.Oid),
		Value: sarama.StringEncoder(e.Oid),
	}
	return mess, nil
}

// Produce event to topic
func Produce(e Event) error {
	mess, err := e.payloadBuild()
	if err != nil {
		return err
	}
	_, _, err = producer.SendMessage(mess)
	if err != nil {
		return ErrMessage
	}
	return nil
}

func newDataCollector() sarama.SyncProducer {
	brokerList := os.Getenv("KAFKA_BROKERS")
	if brokerList == "" {
		return nil
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	syncProducer, err := sarama.NewSyncProducer(strings.Split(brokerList, ","), config)
	if err != nil {
		// log.Fatalln("Failed to start Kafka producer:", err)
		return nil
	}
	return syncProducer
}

func init() {
	producer = newDataCollector()
}
