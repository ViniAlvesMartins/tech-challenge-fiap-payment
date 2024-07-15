package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/pkg/sqs"
	"log"
	"sync"
)

type Handler interface {
	Handle(b []byte) error
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
		fmt.Println("Waiting for message...")
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

		if err = c.handler.Handle([]byte(body.Message)); err != nil {
			log.Println(err.Error())
			continue
		}

		if err = c.service.DeleteMessage(ctx, *m.ReceiptHandle); err != nil {
			log.Println(err.Error())
			continue
		}

		select {
		case <-ctx.Done():
			fmt.Println("Closing consumer...")
			return
		}
	}
}
