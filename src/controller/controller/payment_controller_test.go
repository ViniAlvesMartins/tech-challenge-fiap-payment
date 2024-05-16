package controller

import (
	"bytes"
	"encoding/json"
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

		loggerMock := slog.New(slog.NewTextHandler(os.Stderr, nil))

		orderUseCaseMock := mock.NewMockOrderUseCase(ctrl)
		orderUseCaseMock.EXPECT().GetById(1).Return(order, nil).Times(1)

		paymentUseCaseMock := mock.NewMockPaymentUseCase(ctrl)
		paymentUseCaseMock.EXPECT().CreateQRCode(order).Return(qrCode, nil)

		body := input.PaymentDto{Type: string(enum.QRCODE)}
		jsonBody, _ := json.Marshal(body)
		bodyReader := bytes.NewReader(jsonBody)

		vars := map[string]string{
			"orderId": "1",
		}

		r, _ := http.NewRequest("POST", "/payments", bodyReader)
		w := httptest.NewRecorder()
		r = mux.SetURLVars(r, vars)

		c := NewPaymentController(paymentUseCaseMock, loggerMock, orderUseCaseMock)
		c.CreatePayment(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		//assert.Equal(t, []byte("abcd"), w.Body.Bytes())
	})
}
