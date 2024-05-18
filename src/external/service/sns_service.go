package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/enum"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SnsService struct{}

func NewSnsService() *SnsService { return &SnsService{} }

func NewConnection() (*sns.Client, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),

		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "", SecretAccessKey: "",
			},
		}),
	)

	if err != nil {
		panic(err)
	}

	client := sns.NewFromConfig(cfg)

	return client, nil
}

func (s *SnsService) SendMessage(paymentId int, status enum.PaymentStatus) (bool, error) {

	client, _ := NewConnection()

	message := &Message{
		PaymentId: paymentId,
		Status:    status,
	}

	fmt.Println("teste1234")
	fmt.Println(message)

	messageJs, _ := json.Marshal(message)

	fmt.Println("teste12")
	fmt.Println(messageJs)

	snsMessage := string(messageJs)

	fmt.Println("teste")
	fmt.Println(snsMessage)

	input := &sns.PublishInput{
		Message:  aws.String(snsMessage),
		TopicArn: aws.String("arn:aws:sns:us-east-1:408004708958:payments-sns"),
	}

	result, err := client.Publish(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to publish message, %v", err)
	}

	fmt.Printf("Message ID: %s\n", *result.MessageId)

	return true, nil
}

type Message struct {
	PaymentId int
	Status    enum.PaymentStatus
}
