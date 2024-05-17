package service

import (
	"encoding/json"
	"fmt"
	response_order_service "github.com/ViniAlvesMartins/tech-challenge-fiap/src/application/modules/response/order_service"
	"io"
	"net/http"
	"strconv"
)

type OrderService struct{}

func NewOrderService() *OrderService { return &OrderService{} }

func (o *OrderService) GetById(orderId int) (*response_order_service.GetByIdResp, error) {

	resp, err := http.Get("http://localhost:8080/orders/" + strconv.Itoa(orderId))

	if err != nil {
		fmt.Println("Erro ao fazer a requisição:", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
		return nil, err
	}

	var order response_order_service.GetByIdResp

	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("erro ao decodificar a resposta JSON: %v", err)
	}

	// Imprime a resposta
	fmt.Println(string(body))

	return &order, nil
}
