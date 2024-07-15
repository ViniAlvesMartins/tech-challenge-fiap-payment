package repository

import (
	"context"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-common/uuid"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/contract"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
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

func (p *PaymentRepository) Create(ctx context.Context, payment entity.Payment) error {
	payment.PaymentID = p.uuid.NewString()

	i, err := attributevalue.Marshal(payment)
	if err != nil {
		return err
	}

	items := i.(*types.AttributeValueMemberM).Value
	input := &dynamodb.PutItemInput{
		Item:      items,
		TableName: aws.String(table),
	}

	if _, err = p.db.PutItem(ctx, input); err != nil {
		return err
	}

	return nil
}

func (p *PaymentRepository) GetLastPaymentStatus(ctx context.Context, orderId int) (*entity.Payment, error) {
	var payment *entity.Payment

	out, err := p.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]types.AttributeValue{
			"order_id": &types.AttributeValueMemberN{Value: strconv.Itoa(orderId)},
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
			"order_id": &types.AttributeValueMemberN{Value: strconv.Itoa(orderId)},
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
