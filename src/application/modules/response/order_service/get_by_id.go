package response_order_service

import "github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"

type GetByIdResp struct {
	Error string          `json:"error"`
	Data  *entity.Order `json:"data"`
}
