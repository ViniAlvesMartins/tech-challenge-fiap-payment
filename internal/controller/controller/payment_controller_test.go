package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	contractmock "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/application/contract/mock"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/controller/serializer/input"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
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
		order := &entity.Order{
			ID:          1,
			ClientId:    nil,
			OrderStatus: enum.OrderStatusAwaitingPayment,
			Amount:      123.45,
			CreatedAt:   time.Now(),
		}

		qrCode := &entity.QRCodePayment{QRCode: "qr_data"}
		body := input.PaymentDto{OrderId: order.ID, Type: string(enum.PaymentTypeQRCode)}
		jsonBody, _ := json.Marshal(body)
		bodyReader := bytes.NewReader(jsonBody)

		r, _ := http.NewRequest("POST", "/payments", bodyReader)
		w := httptest.NewRecorder()

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))

		orderUseCaseMock := contractmock.NewOrderUseCase(t)
		orderUseCaseMock.On("GetById", 1).Return(order, nil).Once()

		paymentUseCaseMock := contractmock.NewPaymentUseCase(t)
		paymentUseCaseMock.On("CreateQRCode", r.Context(), order).Return(qrCode, nil)

		c := NewPaymentController(paymentUseCaseMock, loggerMock, orderUseCaseMock)
		c.CreatePayment(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		wrongDTO := struct {
			Data string
		}{
			Data: "wrong data",
		}

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		orderUseCaseMock := contractmock.NewOrderUseCase(t)
		paymentUseCaseMock := contractmock.NewPaymentUseCase(t)

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
		expectedErr := errors.New("error connecting to database")

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		paymentUseCaseMock := contractmock.NewPaymentUseCase(t)

		orderUseCaseMock := contractmock.NewOrderUseCase(t)
		orderUseCaseMock.On("GetById", 1).Return(nil, expectedErr).Once()

		body := input.PaymentDto{OrderId: 1, Type: string(enum.PaymentTypeQRCode)}
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
		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))

		orderUseCaseMock := contractmock.NewOrderUseCase(t)
		orderUseCaseMock.On("GetById", 1).Return(nil, nil).Once()

		paymentUseCaseMock := contractmock.NewPaymentUseCase(t)

		body := input.PaymentDto{OrderId: 1, Type: string(enum.PaymentTypeQRCode)}
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
		order := &entity.Order{
			ID:          1,
			ClientId:    nil,
			OrderStatus: enum.OrderStatusAwaitingPayment,
			Amount:      123.45,
			CreatedAt:   time.Now(),
		}

		body := input.PaymentDto{OrderId: order.ID, Type: string(enum.PaymentTypeQRCode)}
		jsonBody, _ := json.Marshal(body)
		bodyReader := bytes.NewReader(jsonBody)

		r, _ := http.NewRequest("POST", "/payments", bodyReader)
		w := httptest.NewRecorder()

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		orderUseCaseMock := contractmock.NewOrderUseCase(t)
		orderUseCaseMock.On("GetById", 1).Return(order, nil).Once()

		paymentUseCaseMock := contractmock.NewPaymentUseCase(t)
		paymentUseCaseMock.On("CreateQRCode", r.Context(), order).Return(nil, errors.New("error creating qr code"))

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
		order := &entity.Order{
			ID:          1,
			ClientId:    nil,
			OrderStatus: enum.OrderStatusAwaitingPayment,
			Amount:      123.45,
			CreatedAt:   time.Now(),
		}

		body := input.PaymentDto{OrderId: order.ID, Type: string(enum.PaymentTypeQRCode)}
		jsonBody, _ := json.Marshal(body)
		bodyReader := bytes.NewReader(jsonBody)

		r, _ := http.NewRequest("POST", "/payments", bodyReader)
		w := httptest.NewRecorder()

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		orderUseCaseMock := contractmock.NewOrderUseCase(t)
		orderUseCaseMock.On("GetById", 1).Return(order, nil).Once()

		paymentUseCaseMock := contractmock.NewPaymentUseCase(t)
		paymentUseCaseMock.On("CreateQRCode", r.Context(), order).Return(nil, nil)

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
		vars := map[string]string{
			"id": "1",
		}

		r, _ := http.NewRequest("GET", "/payment/{paymentId}/status", nil)
		w := httptest.NewRecorder()
		r = mux.SetURLVars(r, vars)

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		orderUseCaseMock := contractmock.NewOrderUseCase(t)

		paymentUseCaseMock := contractmock.NewPaymentUseCase(t)
		paymentUseCaseMock.On("GetLastPaymentStatus", r.Context(), 1).Return(enum.PaymentStatusPending, nil)

		c := NewPaymentController(paymentUseCaseMock, loggerMock, orderUseCaseMock)
		c.GetLastPaymentStatus(w, r)

		jsonResponse, _ := json.Marshal(Response{
			Error: "",
			Data: GetLastPaymentStatus{
				OrderId:       1,
				PaymentStatus: enum.PaymentStatusPending,
			},
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, string(jsonResponse), string(w.Body.Bytes()))
	})

	t.Run("error converting order id to int", func(t *testing.T) {
		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		orderUseCaseMock := contractmock.NewOrderUseCase(t)
		paymentUseCaseMock := contractmock.NewPaymentUseCase(t)

		vars := map[string]string{
			"id": "not a number",
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
		vars := map[string]string{
			"id": "1",
		}

		r, _ := http.NewRequest("GET", "/payments/{paymentId}/status", nil)
		w := httptest.NewRecorder()
		r = mux.SetURLVars(r, vars)

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		orderUseCaseMock := contractmock.NewOrderUseCase(t)

		paymentUseCaseMock := contractmock.NewPaymentUseCase(t)
		paymentUseCaseMock.On("GetLastPaymentStatus", r.Context(), 1).Return(enum.PaymentStatusPending, errors.New("error getting last payment status"))

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
		vars := map[string]string{
			"id": "1",
		}

		r, _ := http.NewRequest("POST", "/payments/{paymentId}/notification-payments", nil)
		w := httptest.NewRecorder()
		r = mux.SetURLVars(r, vars)

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))
		orderUseCaseMock := contractmock.NewOrderUseCase(t)

		paymentUseCaseMock := contractmock.NewPaymentUseCase(t)
		paymentUseCaseMock.On("ConfirmedPaymentNotification", r.Context(), 1).Return(nil)

		c := NewPaymentController(paymentUseCaseMock, loggerMock, orderUseCaseMock)
		c.Notification(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}
