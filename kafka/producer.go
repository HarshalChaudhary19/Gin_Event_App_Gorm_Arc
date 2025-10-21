package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

const EventsTopic = "events_created"
const UsersTopic = "users_created"

type EventMessage struct {
	// Id          int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Id          int    `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null" json:"name" binding:"required,min=3"`
	OwnerId     int    `gorm:"not null" json:"owner_id" binding:"required" `
	Description string `gorm:"not null" json:"description" binding:"required,min=10"`
	Date        string `gorm:"not null" json:"date" binding:"required"`
	Location    string `gorm:"not null" json:"location" binding:"required,min=3"`
}
type UserSent struct {
	// Id       int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Id       int    `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"unique;not null" json:"email"`
	Name     string `gorm:"not null" json:"name"`
	Password string `json:"-"`
}

var writer *kafka.Writer

func InitProducer(brokerAddr string) {
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{brokerAddr},
		Topic:    UsersTopic,
		Balancer: &kafka.Hash{}, // keep ordering for same key
		Async:    false,         // wait for ack (false => synchronous)
	})
}

// This is for Publishing Events
func PublishEvent(ctx context.Context, key string, ev UserSent) error {
	if writer == nil {
		return fmt.Errorf("kafka writer not initialized")
	}
	payload, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Key:   []byte(key),
		Value: payload,
		Time:  time.Now(),
	}
	return writer.WriteMessages(ctx, msg)
}

func CloseProducer(ctx context.Context) error {
	if writer != nil {
		return writer.Close()
	}
	return nil
}
