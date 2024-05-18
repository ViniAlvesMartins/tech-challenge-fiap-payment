package service

import (
	"encoding/json"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/infra"
	response_order_service "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/application/modules/response/order_service"
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
		//url := fmt.Sprintf("%s/%s", config.OrdersURL, "orders/1")

		response := response_order_service.GetByIdResp{
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
}
