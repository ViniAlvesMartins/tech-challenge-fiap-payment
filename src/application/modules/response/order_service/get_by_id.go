package response_order_service

import "github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"

type GetByIdResp struct {
	error string          `json:"error"`
	data  []*entity.Order `json:"data"`
}