package repository

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log/slog"
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

	//input := &dynamodb.PutItemInput{
	//	Item: map[string]types.AttributeValue{
	//		"id":      &types.AttributeValueMemberS{Value: id},
	//		"orderID": &types.AttributeValueMemberN{Value: uuid.New().String()},
	//		"type":    &types.AttributeValueMemberS{Value: uuid.New().String()},
	//		"status":  &types.AttributeValueMemberS{Value: uuid.New().String()},
	//		"amount":  &types.AttributeValueMemberS{Value: uuid.New().String()},
	//	},
	//	TableName: aws.String(table),
	//}

	//out, err := p.db.PutItem(context.TODO(), input)

	return &payment, nil
}

func (p *PaymentRepository) GetLastPaymentStatus(orderId int) (*entity.Payment, error) {
	/*	var payment entity.Payment

		result := p.db.Order("payments.created_at desc").Where("payments.order_id= ?", orderId).Find(&payment)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return &payment, nil*/
	return nil, nil
}
