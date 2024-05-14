package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"

	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
)

type PaymentRepository struct {
	db  *mongo.Database
	log *slog.Logger
}

func NewPaymentRepository(db *mongo.Database, log *slog.Logger) *PaymentRepository {
	return &PaymentRepository{
		db:  db,
		log: log,
	}
}

func (p *PaymentRepository) Create(payment entity.Payment) (*entity.Payment, error) {

	collection_name := "payments"

	collection := p.db.Collection(collection_name)

	value, err := payment.GetJSONValue()

	if err != nil {
		return nil, errors.New("Error convert domain from JSON")
	}

	result, err := collection.InsertOne(context.Background(), value)

	if err != nil {
		return &payment, errors.New("create payment from repository has failed")
	}

	payment.ID = result.InsertedID.(string)

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
