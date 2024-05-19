package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/contract/mock"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/enum"
	uuidmock "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/pkg/uuid/mock"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"strconv"
	"testing"
)

func TestPaymentRepository_Create(t *testing.T) {
	t.Run("create payment successfully", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		paymentId := uuid.New()

		uuidMock := uuidmock.NewMockInterface(ctrl)
		uuidMock.EXPECT().NewString().Return(paymentId.String()).Times(1)

		db := mock.NewMockDynamoDB(ctrl)
		db.EXPECT().PutItem(ctx, gomock.Any()).Times(1).Return(&dynamodb.PutItemOutput{}, nil)

		p := entity.Payment{
			PaymentID:    paymentId.String(),
			OrderID:      1,
			Type:         enum.QRCODE,
			CurrentState: enum.PENDING,
			Amount:       123.45,
		}

		repo := NewPaymentRepository(db, logger, uuidMock)
		payment, err := repo.Create(ctx, p)

		assert.IsType(t, entity.Payment{}, *payment)
		assert.Nil(t, err)
	})

	t.Run("error creating payment", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("database error")
		paymentId := uuid.New()

		uuidMock := uuidmock.NewMockInterface(ctrl)
		uuidMock.EXPECT().NewString().Return(paymentId.String()).Times(1)

		db := mock.NewMockDynamoDB(ctrl)
		db.EXPECT().PutItem(ctx, gomock.Any()).Times(1).Return(nil, expectedErr)

		p := entity.Payment{
			PaymentID:    paymentId.String(),
			OrderID:      1,
			Type:         enum.QRCODE,
			CurrentState: enum.PENDING,
			Amount:       123.45,
		}

		repo := NewPaymentRepository(db, logger, uuidMock)
		payment, err := repo.Create(ctx, p)

		assert.Nil(t, payment)
		assert.ErrorIs(t, expectedErr, err)
	})
}

func TestPaymentRepository_GetLastPaymentStatus(t *testing.T) {
	t.Run("get last payment status successfully", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		paymentId := uuid.New()

		uuidMock := uuidmock.NewMockInterface(ctrl)

		p := &entity.Payment{
			PaymentID:    paymentId.String(),
			OrderID:      1,
			Type:         enum.QRCODE,
			CurrentState: enum.PENDING,
			Amount:       123.45,
		}

		dynamoOutput := &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"orderId":      &types.AttributeValueMemberN{Value: strconv.Itoa(p.OrderID)},
				"paymentId":    &types.AttributeValueMemberS{Value: p.PaymentID},
				"type":         &types.AttributeValueMemberS{Value: string(p.Type)},
				"currentState": &types.AttributeValueMemberS{Value: string(p.CurrentState)},
				"amount":       &types.AttributeValueMemberN{Value: fmt.Sprint(p.Amount)},
			},
		}

		db := mock.NewMockDynamoDB(ctrl)
		db.EXPECT().GetItem(ctx, gomock.Any()).Times(1).Return(dynamoOutput, nil)

		repo := NewPaymentRepository(db, logger, uuidMock)
		payment, err := repo.GetLastPaymentStatus(ctx, p.OrderID)

		assert.Nil(t, err)
		assert.Equal(t, *p, *payment)
	})

	t.Run("database error getting last payment status", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("database error")
		paymentId := uuid.New()

		uuidMock := uuidmock.NewMockInterface(ctrl)

		p := &entity.Payment{
			PaymentID:    paymentId.String(),
			OrderID:      1,
			Type:         enum.QRCODE,
			CurrentState: enum.PENDING,
			Amount:       123.45,
		}

		db := mock.NewMockDynamoDB(ctrl)
		db.EXPECT().GetItem(ctx, gomock.Any()).Times(1).Return(nil, expectedErr)

		repo := NewPaymentRepository(db, logger, uuidMock)
		payment, err := repo.GetLastPaymentStatus(ctx, p.OrderID)

		assert.ErrorIs(t, expectedErr, err)
		assert.Nil(t, payment)
	})

	t.Run("unmarshall db reponse error", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("unmarshal failed, cannot unmarshal string into Go value type int")
		paymentId := uuid.New()

		uuidMock := uuidmock.NewMockInterface(ctrl)

		p := &entity.Payment{
			PaymentID:    paymentId.String(),
			OrderID:      1,
			Type:         enum.QRCODE,
			CurrentState: enum.PENDING,
			Amount:       123.45,
		}

		dynamoOutput := &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"orderId":      &types.AttributeValueMemberS{Value: strconv.Itoa(p.OrderID)},
				"paymentId":    &types.AttributeValueMemberS{Value: p.PaymentID},
				"type":         &types.AttributeValueMemberS{Value: string(p.Type)},
				"currentState": &types.AttributeValueMemberS{Value: string(p.CurrentState)},
				"amount":       &types.AttributeValueMemberN{Value: fmt.Sprint(p.Amount)},
			},
		}

		db := mock.NewMockDynamoDB(ctrl)
		db.EXPECT().GetItem(ctx, gomock.Any()).Times(1).Return(dynamoOutput, nil)

		repo := NewPaymentRepository(db, logger, uuidMock)
		payment, err := repo.GetLastPaymentStatus(ctx, p.OrderID)

		assert.Errorf(t, expectedErr, err.Error())
		assert.Nil(t, payment)
	})
}

func TestPaymentRepository_UpdateStatus(t *testing.T) {
	t.Run("update status successfully", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

		uuidMock := uuidmock.NewMockInterface(ctrl)

		db := mock.NewMockDynamoDB(ctrl)
		db.EXPECT().UpdateItem(ctx, gomock.Any()).Times(1).Return(&dynamodb.UpdateItemOutput{}, nil)

		repo := NewPaymentRepository(db, logger, uuidMock)
		err := repo.UpdateStatus(ctx, 1, enum.CONFIRMED)

		assert.Nil(t, err)
	})

	t.Run("error updating status", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("database error")

		uuidMock := uuidmock.NewMockInterface(ctrl)

		db := mock.NewMockDynamoDB(ctrl)
		db.EXPECT().UpdateItem(ctx, gomock.Any()).Times(1).Return(nil, expectedErr)

		repo := NewPaymentRepository(db, logger, uuidMock)
		err := repo.UpdateStatus(ctx, 1, enum.CONFIRMED)

		assert.ErrorIs(t, expectedErr, err)
	})
}
