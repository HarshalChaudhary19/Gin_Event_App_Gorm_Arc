package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func StartConsumerGroup(ctx context.Context, brokerAddr, groupID string, handler func(EventMessage) error) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokerAddr},
		GroupID:     groupID,
		Topic:       EventsTopic,
		MinBytes:    10e3,
		MaxBytes:    10e6,
		StartOffset: kafka.LastOffset,
	})
	go func() {
		defer r.Close()
		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				// If context canceled, exit gracefully
				if ctx.Err() != nil {
					return
				}
				log.Printf("kafka read error: %v", err)
				time.Sleep(time.Second)
				continue
			}
			var ev EventMessage
			if err := json.Unmarshal(m.Value, &ev); err != nil {
				log.Printf("unmarshal error: %v", err)
				continue
			}
			if err := handler(ev); err != nil {
				// Handler-level error handling: log, maybe send to DLQ later
				log.Printf("handler error: %v", err)
			}
		}
	}()
	return nil
}
