package sqs

import (
	"context"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/pkg/sqs"
)

type Handler interface {
	Handle() error
}

type Consumer struct {
	service *sqs.Service
	handler Handler
}

func NewConsumer(s *sqs.Service) *Consumer {
	return &Consumer{service: s}
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		message, err := c.service.ReceiveMessage(ctx)
		if err != nil {
			return err
		}

		if message == nil {
			continue
		}

	}
}
