package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/infra"
	responseorderservice "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/modules/response/order_service"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/enum"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestOrderService_GetById(t *testing.T) {
	t.Run("get order successfully", func(t *testing.T) {
		config, _ := infra.NewConfig()
		client := resty.New().SetBaseURL(config.OrdersURL)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

		response := responseorderservice.GetByIdResp{
			Error: "",
			Data: &entity.Order{
				ID:          1,
				ClientId:    nil,
				StatusOrder: enum.AWAITING_PAYMENT,
				Amount:      123.45,
				CreatedAt:   time.Now(),
			},
		}

		httpmock.ActivateNonDefault(client.GetClient())
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", fmt.Sprintf("%s%s", config.OrdersURL, "/orders/1"),
			httpmock.NewJsonResponderOrPanic(200, response))

		orderService := NewOrderService(client, logger)
		order, err := orderService.GetById(1)

		expectedJsonResponse, _ := json.Marshal(response)
		jsonResponse, _ := json.Marshal(order)

		assert.Equal(t, expectedJsonResponse, jsonResponse)
		assert.Nil(t, err)
	})

	t.Run("error making request", func(t *testing.T) {
		config, _ := infra.NewConfig()
		client := resty.New().SetBaseURL("http://nonexistent.app.com")
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("Get \"orders\": no responder found")

		response := responseorderservice.GetByIdResp{
			Error: "",
			Data: &entity.Order{
				ID:          1,
				ClientId:    nil,
				StatusOrder: enum.AWAITING_PAYMENT,
				Amount:      123.45,
				CreatedAt:   time.Now(),
			},
		}

		httpmock.ActivateNonDefault(client.GetClient())
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", fmt.Sprintf("%s%s", config.OrdersURL, "/teste./1"),
			httpmock.NewJsonResponderOrPanic(200, response))

		orderService := NewOrderService(client, logger)
		order, err := orderService.GetById(1)

		assert.Errorf(t, expectedErr, err.Error())
		assert.Nil(t, order)
	})

	t.Run("request with error code response", func(t *testing.T) {
		config, _ := infra.NewConfig()
		client := resty.New().SetBaseURL(config.OrdersURL)
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
		expectedErr := errors.New("request response error")

		response := responseorderservice.GetByIdResp{
			Error: "Internal server error",
			Data:  nil,
		}

		httpmock.ActivateNonDefault(client.GetClient())
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", fmt.Sprintf("%s%s", config.OrdersURL, "/orders/1"),
			httpmock.NewJsonResponderOrPanic(500, response))

		orderService := NewOrderService(client, logger)
		order, err := orderService.GetById(1)

		assert.Errorf(t, expectedErr, err.Error())
		assert.Nil(t, order)
	})
}
