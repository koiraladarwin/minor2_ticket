package kafka

import (
	"context"
	"encoding/json"
	"os"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer() *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(os.Getenv("KAFKA_BROKER")),
			Balancer: &kafka.LeastBytes{},
		},
		topic: os.Getenv("KAFKA_TOPIC"),
	}
}

func (p *Producer) Publish(
	ctx context.Context,
	topic string,
	event any,
) error {

	b, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(
		ctx,
		kafka.Message{
			Topic: topic,
			Value: b,
		},
	)
}
