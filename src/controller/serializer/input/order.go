package input

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"
)

type OrderDto struct {
	ClientId *int `json:"client_id"`
	Products []struct {
		ID int `json:"id"`
	} `json:"products"`
}

func (o OrderDto) ConvertToEntity() entity.Order {
	var products []*entity.Product

	for _, p := range o.Products {
		products = append(products, &entity.Product{ID: p.ID})
	}

	return entity.Order{
		ClientId: o.ClientId,
		Products: products,
	}
}
