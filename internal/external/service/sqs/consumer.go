package sqs

import (
	"context"
	"encoding/json"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/pkg/sqs"
	"log"
)

type Handler interface {
	Handle(a any) error
}

type Consumer[T interface{}] struct {
	service *sqs.Service
	handler Handler
}

type MessageBody struct {
	Type      string
	MessageId string
	TopicArn  string
	Message   string
	Timestamp string
}

func NewConsumer(s *sqs.Service) *Consumer[interface{}] {
	return &Consumer[interface{}]{service: s}
}

func (c *Consumer[T]) Start(ctx context.Context) error {
	for {
		m, err := c.service.ReceiveMessage(ctx)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if m == nil {
			continue
		}

		var body *MessageBody
		if err = json.Unmarshal([]byte(*m.Body), &body); err != nil {
			log.Println(err.Error())
			continue
		}

		var message T
		if err = json.Unmarshal([]byte(body.Message), &message); err != nil {
			log.Println(err.Error())
			continue
		}

		if err = c.handler.Handle(message); err != nil {
			log.Println(err.Error())
			continue
		}
	}
}
