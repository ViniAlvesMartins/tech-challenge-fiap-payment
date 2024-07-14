package repository

import (
	"context"
	"errors"
	"fmt"
	contractmock "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/contract/mock"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	uuidmock "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/pkg/uuid/mock"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"os"
	"strconv"
	"testing"
)

func TestPaymentRepository_Create(t *testing.T) {
	t.Run("create payment successfully", func(t *testing.T) {
		ctx := context.Background()
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		paymentId := uuid.New()

		uuidMock := uuidmock.NewInterface(t)
		uuidMock.On("NewString").Return(paymentId.String()).Once()

		db := contractmock.NewDynamoDB(t)
		db.On("PutItem", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).Once().Return(&dynamodb.PutItemOutput{}, nil)

		p := entity.Payment{
			PaymentID:    paymentId.String(),
			OrderID:      1,
			Type:         enum.QRCODE,
			CurrentState: enum.PENDING,
			Amount:       123.45,
		}

		repo := NewPaymentRepository(db, logger, uuidMock)
		err := repo.Create(ctx, p)

		assert.Nil(t, err)
	})

	t.Run("error creating payment", func(t *testing.T) {
		ctx := context.Background()
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("database error")
		paymentId := uuid.New()

		uuidMock := uuidmock.NewInterface(t)
		uuidMock.On("NewString").Return(paymentId.String()).Once()

		db := contractmock.NewDynamoDB(t)
		db.On("PutItem", ctx, mock.AnythingOfType("*dynamodb.PutItemInput")).Once().Return(nil, expectedErr)

		p := entity.Payment{
			PaymentID:    paymentId.String(),
			OrderID:      1,
			Type:         enum.QRCODE,
			CurrentState: enum.PENDING,
			Amount:       123.45,
		}

		repo := NewPaymentRepository(db, logger, uuidMock)
		err := repo.Create(ctx, p)

		assert.ErrorIs(t, expectedErr, err)
	})
}

func TestPaymentRepository_GetLastPaymentStatus(t *testing.T) {
	t.Run("get last payment status successfully", func(t *testing.T) {
		ctx := context.Background()
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		paymentId := uuid.New()

		uuidMock := uuidmock.NewInterface(t)

		p := &entity.Payment{
			PaymentID:    paymentId.String(),
			OrderID:      1,
			Type:         enum.QRCODE,
			CurrentState: enum.PENDING,
			Amount:       123.45,
		}

		dynamoOutput := &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"order_id":   &types.AttributeValueMemberN{Value: strconv.Itoa(p.OrderID)},
				"payment_id": &types.AttributeValueMemberS{Value: p.PaymentID},
				"type":       &types.AttributeValueMemberS{Value: string(p.Type)},
				"status":     &types.AttributeValueMemberS{Value: string(p.CurrentState)},
				"amount":     &types.AttributeValueMemberN{Value: fmt.Sprint(p.Amount)},
			},
		}

		db := contractmock.NewDynamoDB(t)
		db.On("GetItem", ctx, mock.AnythingOfType("*dynamodb.GetItemInput")).Once().Return(dynamoOutput, nil)

		repo := NewPaymentRepository(db, logger, uuidMock)
		payment, err := repo.GetLastPaymentStatus(ctx, p.OrderID)

		assert.Nil(t, err)
		assert.Equal(t, *p, *payment)
	})

	t.Run("database error getting last payment status", func(t *testing.T) {
		ctx := context.Background()
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("database error")
		paymentId := uuid.New()

		uuidMock := uuidmock.NewInterface(t)

		p := &entity.Payment{
			PaymentID:    paymentId.String(),
			OrderID:      1,
			Type:         enum.QRCODE,
			CurrentState: enum.PENDING,
			Amount:       123.45,
		}

		db := contractmock.NewDynamoDB(t)
		db.On("GetItem", ctx, mock.AnythingOfType("*dynamodb.GetItemInput")).Once().Return(nil, expectedErr)

		repo := NewPaymentRepository(db, logger, uuidMock)
		payment, err := repo.GetLastPaymentStatus(ctx, p.OrderID)

		assert.ErrorIs(t, expectedErr, err)
		assert.Nil(t, payment)
	})

	t.Run("unmarshall db response error", func(t *testing.T) {
		ctx := context.Background()
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("unmarshal failed, cannot unmarshal string into Go value type float32")
		paymentId := uuid.New()

		uuidMock := uuidmock.NewInterface(t)

		p := &entity.Payment{
			PaymentID:    paymentId.String(),
			OrderID:      1,
			Type:         enum.QRCODE,
			CurrentState: enum.PENDING,
			Amount:       123.45,
		}

		dynamoOutput := &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"order_id":   &types.AttributeValueMemberN{Value: fmt.Sprint(p.OrderID)},
				"payment_id": &types.AttributeValueMemberS{Value: p.PaymentID},
				"type":       &types.AttributeValueMemberS{Value: string(p.Type)},
				"status":     &types.AttributeValueMemberS{Value: string(p.CurrentState)},
				"amount":     &types.AttributeValueMemberS{Value: fmt.Sprint(p.Amount)},
			},
		}

		db := contractmock.NewDynamoDB(t)
		db.On("GetItem", ctx, mock.AnythingOfType("*dynamodb.GetItemInput")).Once().Return(dynamoOutput, nil)

		repo := NewPaymentRepository(db, logger, uuidMock)
		payment, err := repo.GetLastPaymentStatus(ctx, p.OrderID)

		assert.Error(t, expectedErr, err)
		assert.Nil(t, payment)
	})
}

func TestPaymentRepository_UpdateStatus(t *testing.T) {
	t.Run("update status successfully", func(t *testing.T) {
		ctx := context.Background()
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

		uuidMock := uuidmock.NewInterface(t)

		db := contractmock.NewDynamoDB(t)
		db.On("UpdateItem", ctx, mock.AnythingOfType("*dynamodb.UpdateItemInput")).Once().Return(&dynamodb.UpdateItemOutput{}, nil)

		repo := NewPaymentRepository(db, logger, uuidMock)
		err := repo.UpdateStatus(ctx, 1, enum.CONFIRMED)

		assert.Nil(t, err)
	})

	t.Run("error updating status", func(t *testing.T) {
		ctx := context.Background()
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("database error")

		uuidMock := uuidmock.NewInterface(t)

		db := contractmock.NewDynamoDB(t)
		db.On("UpdateItem", ctx, mock.AnythingOfType("*dynamodb.UpdateItemInput")).Once().Return(nil, expectedErr)

		repo := NewPaymentRepository(db, logger, uuidMock)
		err := repo.UpdateStatus(ctx, 1, enum.CONFIRMED)

		assert.ErrorIs(t, expectedErr, err)
	})
}
