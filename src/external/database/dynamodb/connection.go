package dynamodb

import (
	"context"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/infra"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewConnection(conf infra.Config) (*dynamodb.Client, error) {

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

	client := dynamodb.NewFromConfig(cfg)

	return client, nil

}
