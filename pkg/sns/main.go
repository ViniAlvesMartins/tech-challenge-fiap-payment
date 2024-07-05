package sns

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type Service struct {
	client *sns.Client
	topic  string
}

func NewConnection(ctx context.Context, topic string) (*Service, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &Service{
		client: sns.NewFromConfig(cfg),
		topic:  topic,
	}, nil
}

func (s *Service) Publish(ctx context.Context, data any) error {
	message, err := json.Marshal(data)
	if err != nil {
		return err
	}

	input := &sns.PublishInput{
		Message:  aws.String(string(message)),
		TopicArn: aws.String(s.topic),
	}

	_, err = s.client.Publish(ctx, input)
	return err
}
