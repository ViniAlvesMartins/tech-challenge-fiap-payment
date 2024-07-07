package repository

import (
	"context"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/pkg/uuid"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"

	"log/slog"
	"strconv"
)

const table = "payments"

type PaymentRepository struct {
	db   contract.DynamoDB
	log  *slog.Logger
	uuid uuid.Interface
}

func NewPaymentRepository(db contract.DynamoDB, log *slog.Logger, u uuid.Interface) *PaymentRepository {
	return &PaymentRepository{
		db:   db,
		log:  log,
		uuid: u,
	}
}

func (p *PaymentRepository) Create(ctx context.Context, payment entity.Payment) (*entity.Payment, error) {
	id := p.uuid.NewString()

	input := &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"order_id":      &types.AttributeValueMemberS{Value: payment.OrderID},
			"payment_id":    &types.AttributeValueMemberS{Value: id},
			"type":          &types.AttributeValueMemberS{Value: string(payment.Type)},
			"current_state": &types.AttributeValueMemberS{Value: string(payment.CurrentState)},
			"amount":        &types.AttributeValueMemberN{Value: fmt.Sprint(payment.Amount)},
		},
		TableName: aws.String(table),
	}

	if _, err := p.db.PutItem(ctx, input); err != nil {
		return nil, err
	}

	payment.PaymentID = id
	return &payment, nil
}

func (p *PaymentRepository) GetLastPaymentStatus(ctx context.Context, orderId int) (*entity.Payment, error) {
	payment := &entity.Payment{}

	out, err := p.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]types.AttributeValue{
			"order_id": &types.AttributeValueMemberS{Value: strconv.Itoa(orderId)},
		},
	})

	if err != nil {
		return nil, err
	}

	if err = attributevalue.UnmarshalMap(out.Item, &payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (p *PaymentRepository) UpdateStatus(ctx context.Context, orderId int, status enum.PaymentStatus) error {
	if _, err := p.db.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(table),
		Key: map[string]types.AttributeValue{
			"order_id": &types.AttributeValueMemberS{Value: strconv.Itoa(orderId)},
		},
		UpdateExpression: aws.String("set current_state = :current_state"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":current_state": &types.AttributeValueMemberS{Value: string(status)},
		},
	}); err != nil {
		return err
	}

	return nil
}
