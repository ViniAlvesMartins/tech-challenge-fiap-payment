package mongodb

import (
	"context"
	"fmt"
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
				AccessKeyID: "local", SecretAccessKey: "local", SessionToken: "",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)

	if err != nil {
		panic(err)
	}

	client := dynamodb.NewFromConfig(cfg)

	client.Options()

	input := &dynamodb.DescribeTableInput{
		TableName: aws.String("payments"),
	}

	resp, erro := client.DescribeTable(context.TODO(), input)

	if erro != nil {
		fmt.Println("Got error calling DescribeTable:")
		fmt.Println(erro)
		fmt.Println(erro.Error())
	} else {
		fmt.Println("Successfully list table to table")
		fmt.Println(resp)
		fmt.Println(erro)
	}

	return client, nil

}
