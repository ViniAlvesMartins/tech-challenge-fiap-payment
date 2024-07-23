package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-common/sqs"
	"log"
	"sync"
)

type Handler interface {
	Handle(ctx context.Context, b []byte) error
}

type Consumer struct {
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

func NewConsumer(s *sqs.Service, h Handler) *Consumer {
	return &Consumer{
		service: s,
		handler: h,
	}
}

func (c *Consumer) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Closing consumer...")
			return
		default:
		}

		c.consume(ctx)
	}
}

func (c *Consumer) consume(ctx context.Context) {
	fmt.Println("Waiting for message...")
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
		fmt.Println(err.Error())
		return
	}

	if err = c.handler.Handle(ctx, []byte(body.Message)); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = c.service.DeleteMessage(ctx, *m.ReceiptHandle); err != nil {
		fmt.Println(err.Error())
		return
	}
}
