//go:generate mockery --all --output mock --outpkg mock
package contract

import (
	"context"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment entity.Payment) error
	GetLastPaymentStatus(ctx context.Context, orderId int) (*entity.Payment, error)
	UpdateStatus(ctx context.Context, orderId int, status enum.PaymentStatus) error
}

type DynamoDB interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
}
