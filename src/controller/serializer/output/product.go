package output

import "github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"

type ProductDto struct {
	ID          int     `json:"id"`
	NameProduct string  `json:"name_product"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	CategoryId  int     `json:"category_id"`
	Active      bool    `json:"active"`
}

func ProductFromEntity(product entity.Product) ProductDto {
	return ProductDto{
		ID:          product.ID,
		NameProduct: product.NameProduct,
		Description: product.Description,
		Price:       product.Price,
		CategoryId:  product.CategoryId,
		Active:      product.Active,
	}
}

func ProductListFromEntity(products []entity.Product) []ProductDto {
	var productsDto []ProductDto
	for _, p := range products {
		productsDto = append(productsDto, ProductFromEntity(p))
	}

	return productsDto
}
