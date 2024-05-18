package repository

import (
	"context"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/enum"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"log/slog"
	"strconv"
)

type PaymentRepository struct {
	db  *dynamodb.Client
	log *slog.Logger
}

func NewPaymentRepository(db *dynamodb.Client, log *slog.Logger) *PaymentRepository {
	return &PaymentRepository{
		db:  db,
		log: log,
	}
}

type Item struct {
	Id      string `json:"pk"`
	OrderId string `json:"sk"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Amount  string `json:"amount"`
}

func (p *PaymentRepository) Create(payment entity.Payment) (*entity.Payment, error) {

	table := "payments"
	id := uuid.New().String()

	input := &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"orderId":      &types.AttributeValueMemberS{Value: strconv.Itoa(payment.OrderID)},
			"paymentId":    &types.AttributeValueMemberS{Value: id},
			"type":         &types.AttributeValueMemberS{Value: string(payment.Type)},
			"currentState": &types.AttributeValueMemberS{Value: string(payment.CurrentState)},
			"amount":       &types.AttributeValueMemberN{Value: fmt.Sprint(payment.Amount)},
		},
		TableName: aws.String(table),
	}

	_, err := p.db.PutItem(context.TODO(), input)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
	}

	payment.PaymentID = id

	return &payment, nil
}

func (p *PaymentRepository) GetLastPaymentStatus(orderId int) (*entity.Payment, error) {

	payment := &entity.Payment{}

	table := "payments"

	out, err := p.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]types.AttributeValue{
			"orderId": &types.AttributeValueMemberS{Value: strconv.Itoa(payment.OrderID)},
		},
	})

	if err != nil {
		panic(err)
	}

	err = attributevalue.UnmarshalMap(out.Item, &payment)

	if err != nil {
		panic(err)
	}

	return payment, nil
}

func (p *PaymentRepository) UpdateStatus(orderId int, status enum.PaymentStatus) (bool, error) {
	table := "payments"

	_, err := p.db.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(table),
		Key: map[string]types.AttributeValue{
			"orderId": &types.AttributeValueMemberS{Value: strconv.Itoa(orderId)},
		},
		UpdateExpression: aws.String("set currentState = :currentState"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":currentState": &types.AttributeValueMemberS{Value: string(status)},
		},
	})

	if err != nil {
		panic(err)
	}

	return true, nil
}
