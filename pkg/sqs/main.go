package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Service struct {
	queueUrl          string
	client            *sqs.Client
	messagesMaxNumber int32
	waitTime          int32
}

func NewConnection(ctx context.Context, region string, queueUrl string, messagesMaxNumber, waitTime int32) (*Service, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	return &Service{
		queueUrl:          queueUrl,
		client:            sqs.NewFromConfig(cfg),
		messagesMaxNumber: messagesMaxNumber,
		waitTime:          waitTime,
	}, nil
}

func (s *Service) ReceiveMessage(ctx context.Context) (*types.Message, error) {
	payload, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &s.queueUrl,
		MaxNumberOfMessages: s.messagesMaxNumber,
		WaitTimeSeconds:     s.waitTime,
	})
	if err != nil {
		return nil, err
	}

	if len(payload.Messages) >= 1 {
		return &payload.Messages[0], nil
	}

	return nil, nil
}

func (s *Service) DeleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &s.queueUrl,
		ReceiptHandle: &receiptHandle,
	})

	return err
}

func (s *Service) SendMessage(ctx context.Context, message string, messageGroupId string) error {
	_, err := s.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:               &s.queueUrl,
		MessageBody:            &message,
		MessageGroupId:         &messageGroupId,
		MessageDeduplicationId: &message,
	})

	return err
}
