package sqs

import (
	"context"
	"encoding/json"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-common/sqs"
	"log"
	"log/slog"
	"sync"
)

type Handler interface {
	Handle(ctx context.Context, b []byte) error
}

type Consumer struct {
	service *sqs.Service
	handler Handler
	logger  *slog.Logger
}

type MessageBody struct {
	Type      string
	MessageId string
	TopicArn  string
	Message   string
	Timestamp string
}

func NewConsumer(s *sqs.Service, h Handler, l *slog.Logger) *Consumer {
	return &Consumer{
		service: s,
		handler: h,
		logger:  l,
	}
}

func (c *Consumer) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Closing consumer...")
			return
		default:
		}

		c.consume(ctx)
	}
}

func (c *Consumer) consume(ctx context.Context) {
	c.logger.Info("Waiting for message...")
	m, err := c.service.ReceiveMessage(ctx)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if m == nil {
		return
	}

	var body *MessageBody
	if err = json.Unmarshal([]byte(*m.Body), &body); err != nil {
		c.logger.Error("error unmarshalling message", slog.String("error", err.Error()))
		return
	}

	if err = c.handler.Handle(ctx, []byte(body.Message)); err != nil {
		c.logger.Error("error handling message", slog.String("error", err.Error()))
		return
	}

	if err = c.service.DeleteMessage(ctx, *m.ReceiptHandle); err != nil {
		c.logger.Error("error deleting message", slog.String("error", err.Error()))
		return
	}
}
