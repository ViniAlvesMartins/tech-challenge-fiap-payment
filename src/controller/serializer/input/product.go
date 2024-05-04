package input

import "github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"

type ProductDto struct {
	NameProduct string  `json:"name_product" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float32 `json:"price" validate:"required,gt=0" error:"O Número deve ser maior que zero"`
	CategoryId  int     `json:"category_id" validate:"required"`
}

func (p *ProductDto) ConvertToEntity() entity.Product {
	return entity.Product{
		NameProduct: p.NameProduct,
		Description: p.Description,
		Price:       p.Price,
		CategoryId:  p.CategoryId,
	}
}
