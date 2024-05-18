package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"log"
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

func (s *SnsService) sendMessage(payment entity.Payment) (bool, error) {

	client, _ := NewConnection()

	js, _ := json.Marshal(payment)
	message := string(js)

	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String("arn:aws:sns:us-east-1:408004708958:payments-sns"),
	}

	result, err := client.Publish(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to publish message, %v", err)
	}

	// Exibir o resultado
	fmt.Printf("Message ID: %s\n", *result.MessageId)

	return true, nil
}
