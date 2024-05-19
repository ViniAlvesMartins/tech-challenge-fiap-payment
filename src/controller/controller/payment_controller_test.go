package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/contract/mock"
	responsepaymentservice "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/modules/response/payment_service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/controller/serializer/input"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/enum"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestPaymentController_CreatePayment(t *testing.T) {
	t.Run("create payment successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		order := &entity.Order{
			ID:          1,
			ClientId:    nil,
			StatusOrder: enum.AWAITING_PAYMENT,
			Amount:      123.45,
			CreatedAt:   time.Now(),
		}

		qrCode := &responsepaymentservice.CreateQRCode{QrData: "qr_data"}
		body := input.PaymentDto{OrderId: order.ID, Type: string(enum.QRCODE)}
		jsonBody, _ := json.Marshal(body)
		bodyReader := bytes.NewReader(jsonBody)

		r, _ := http.NewRequest("POST", "/payments", bodyReader)
		w := httptest.NewRecorder()

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))

		orderUseCaseMock := mock.NewMockOrderUseCase(ctrl)
		orderUseCaseMock.EXPECT().GetById(1).Return(order, nil).Times(1)

		paymentUseCaseMock := mock.NewMockPaymentUseCase(ctrl)
		paymentUseCaseMock.EXPECT().CreateQRCode(r.Context(), order).Return(qrCode, nil)

		c := NewPaymentController(paymentUseCaseMock, loggerMock, orderUseCaseMock)
		c.CreatePayment(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		wrongDTO := struct {
			Data string
		}{
			Data: "wrong data",
		}

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		orderUseCaseMock := mock.NewMockOrderUseCase(ctrl)
		paymentUseCaseMock := mock.NewMockPaymentUseCase(ctrl)

		jsonBody, _ := json.Marshal(wrongDTO)
		bodyReader := bytes.NewReader(jsonBody)

		r, _ := http.NewRequest("POST", "/payments", bodyReader)
		w := httptest.NewRecorder()

		c := NewPaymentController(paymentUseCaseMock, loggerMock, orderUseCaseMock)
		c.CreatePayment(w, r)

		jsonResponse, _ := json.Marshal(Response{
			Error: "Make sure all required fields are sent correctly",
			Data:  nil,
		})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, string(jsonResponse), string(w.Body.Bytes()))
	})

	t.Run("error getting order by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		expectedErr := errors.New("error connecting to database")

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		paymentUseCaseMock := mock.NewMockPaymentUseCase(ctrl)

		orderUseCaseMock := mock.NewMockOrderUseCase(ctrl)
		orderUseCaseMock.EXPECT().GetById(1).Return(nil, expectedErr).Times(1)

		body := input.PaymentDto{OrderId: 1, Type: string(enum.QRCODE)}
		jsonBody, _ := json.Marshal(body)
		bodyReader := bytes.NewReader(jsonBody)

		r, _ := http.NewRequest("POST", "/payments", bodyReader)
		w := httptest.NewRecorder()

		c := NewPaymentController(paymentUseCaseMock, loggerMock, orderUseCaseMock)
		c.CreatePayment(w, r)

		jsonResponse, _ := json.Marshal(Response{
			Error: "Error getting order details",
			Data:  nil,
		})

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, string(jsonResponse), string(w.Body.Bytes()))
	})

	t.Run("order not found error", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))

		orderUseCaseMock := mock.NewMockOrderUseCase(ctrl)
		orderUseCaseMock.EXPECT().GetById(1).Return(nil, nil).Times(1)

		paymentUseCaseMock := mock.NewMockPaymentUseCase(ctrl)

		body := input.PaymentDto{OrderId: 1, Type: string(enum.QRCODE)}
		jsonBody, _ := json.Marshal(body)
		bodyReader := bytes.NewReader(jsonBody)

		r, _ := http.NewRequest("POST", "/payments", bodyReader)
		w := httptest.NewRecorder()

		c := NewPaymentController(paymentUseCaseMock, loggerMock, orderUseCaseMock)
		c.CreatePayment(w, r)

		jsonResponse, _ := json.Marshal(Response{
			Error: "Order not found",
			Data:  nil,
		})

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, string(jsonResponse), string(w.Body.Bytes()))
	})

	t.Run("error creating qr code", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		order := &entity.Order{
			ID:          1,
			ClientId:    nil,
			StatusOrder: enum.AWAITING_PAYMENT,
			Amount:      123.45,
			CreatedAt:   time.Now(),
		}

		body := input.PaymentDto{OrderId: order.ID, Type: string(enum.QRCODE)}
		jsonBody, _ := json.Marshal(body)
		bodyReader := bytes.NewReader(jsonBody)

		r, _ := http.NewRequest("POST", "/payments", bodyReader)
		w := httptest.NewRecorder()

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		orderUseCaseMock := mock.NewMockOrderUseCase(ctrl)
		orderUseCaseMock.EXPECT().GetById(1).Return(order, nil).Times(1)

		paymentUseCaseMock := mock.NewMockPaymentUseCase(ctrl)
		paymentUseCaseMock.EXPECT().CreateQRCode(r.Context(), order).Return(nil, errors.New("error creating qr code"))

		jsonResponse, _ := json.Marshal(Response{
			Error: "Error creating qr code",
			Data:  nil,
		})

		c := NewPaymentController(paymentUseCaseMock, loggerMock, orderUseCaseMock)
		c.CreatePayment(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, string(jsonResponse), string(w.Body.Bytes()))
	})

	t.Run("empty qr code", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		order := &entity.Order{
			ID:          1,
			ClientId:    nil,
			StatusOrder: enum.AWAITING_PAYMENT,
			Amount:      123.45,
			CreatedAt:   time.Now(),
		}

		body := input.PaymentDto{OrderId: order.ID, Type: string(enum.QRCODE)}
		jsonBody, _ := json.Marshal(body)
		bodyReader := bytes.NewReader(jsonBody)

		r, _ := http.NewRequest("POST", "/payments", bodyReader)
		w := httptest.NewRecorder()

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		orderUseCaseMock := mock.NewMockOrderUseCase(ctrl)
		orderUseCaseMock.EXPECT().GetById(1).Return(order, nil).Times(1)

		paymentUseCaseMock := mock.NewMockPaymentUseCase(ctrl)
		paymentUseCaseMock.EXPECT().CreateQRCode(r.Context(), order).Return(nil, nil)

		c := NewPaymentController(paymentUseCaseMock, loggerMock, orderUseCaseMock)
		c.CreatePayment(w, r)

		jsonResponse, _ := json.Marshal(Response{
			Error: "O pagamento para o pedido já foi efetuado",
			Data:  nil,
		})

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, string(jsonResponse), string(w.Body.Bytes()))
	})
}

func TestPaymentController_GetLastPaymentStatus(t *testing.T) {
	t.Run("get last payment status successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		vars := map[string]string{
			"paymentId": "1",
		}

		r, _ := http.NewRequest("GET", "/payment/{paymentId}/status", nil)
		w := httptest.NewRecorder()
		r = mux.SetURLVars(r, vars)

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		orderUseCaseMock := mock.NewMockOrderUseCase(ctrl)

		paymentUseCaseMock := mock.NewMockPaymentUseCase(ctrl)
		paymentUseCaseMock.EXPECT().GetLastPaymentStatus(r.Context(), 1).Return(enum.PENDING, nil)

		c := NewPaymentController(paymentUseCaseMock, loggerMock, orderUseCaseMock)
		c.GetLastPaymentStatus(w, r)

		jsonResponse, _ := json.Marshal(Response{
			Error: "",
			Data: GetLastPaymentStatus{
				OrderId:       1,
				PaymentStatus: enum.PENDING,
			},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, string(jsonResponse), string(w.Body.Bytes()))
	})

	t.Run("error converting order id to int", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		orderUseCaseMock := mock.NewMockOrderUseCase(ctrl)
		paymentUseCaseMock := mock.NewMockPaymentUseCase(ctrl)

		vars := map[string]string{
			"paymentId": "not a number",
		}

		r, _ := http.NewRequest("GET", "/status-payment", nil)
		w := httptest.NewRecorder()
		r = mux.SetURLVars(r, vars)

		c := NewPaymentController(paymentUseCaseMock, loggerMock, orderUseCaseMock)
		c.GetLastPaymentStatus(w, r)

		jsonResponse, _ := json.Marshal(Response{
			Error: "Order id must be an integer",
			Data:  nil,
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, string(jsonResponse), string(w.Body.Bytes()))
	})

	t.Run("error getting last payment status", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		vars := map[string]string{
			"paymentId": "1",
		}

		r, _ := http.NewRequest("GET", "/payments/{paymentId}/status", nil)
		w := httptest.NewRecorder()
		r = mux.SetURLVars(r, vars)

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		orderUseCaseMock := mock.NewMockOrderUseCase(ctrl)

		paymentUseCaseMock := mock.NewMockPaymentUseCase(ctrl)
		paymentUseCaseMock.EXPECT().GetLastPaymentStatus(r.Context(), 1).Return(enum.PENDING, errors.New("error getting last payment status"))

		c := NewPaymentController(paymentUseCaseMock, loggerMock, orderUseCaseMock)
		c.GetLastPaymentStatus(w, r)

		jsonResponse, _ := json.Marshal(Response{
			Error: "Error getting last payment status",
			Data:  nil,
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, string(jsonResponse), string(w.Body.Bytes()))
	})
}

func TestPaymentController_Notification(t *testing.T) {
	t.Run("get last payment status successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		vars := map[string]string{
			"paymentId": "1",
		}

		r, _ := http.NewRequest("POST", "/payments/{paymentId}/notification-payments", nil)
		w := httptest.NewRecorder()
		r = mux.SetURLVars(r, vars)

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		orderUseCaseMock := mock.NewMockOrderUseCase(ctrl)

		paymentUseCaseMock := mock.NewMockPaymentUseCase(ctrl)
		paymentUseCaseMock.EXPECT().PaymentNotification(r.Context(), 1).Return(nil)

		c := NewPaymentController(paymentUseCaseMock, loggerMock, orderUseCaseMock)
		c.Notification(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}
