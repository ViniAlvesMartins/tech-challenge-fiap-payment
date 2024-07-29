package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/enum"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SnsService struct{}

type Message struct {
	OrderId int
	Status  enum.PaymentStatus
}

func NewSnsService() *SnsService { return &SnsService{} }

func NewConnection() (*sns.Client, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)

	if err != nil {
		panic(err)
	}

	client := sns.NewFromConfig(cfg)

	return client, nil
}

func (s *SnsService) SendMessage(paymentId int, status enum.PaymentStatus) error {
	client, _ := NewConnection()

	message := &Message{
		OrderId: paymentId,
		Status:  status,
	}

	messageJs, _ := json.Marshal(message)

	snsMessage := string(messageJs)

	input := &sns.PublishInput{
		Message:  aws.String(snsMessage),
		TopicArn: aws.String("arn:aws:sns:us-east-1:961381910932:update_order_status-topic"),
	}

	result, err := client.Publish(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to publish message, %v", err)
	}

	fmt.Printf("Message ID: %s\n", *result.MessageId)

	return nil
}
