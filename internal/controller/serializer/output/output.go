package output

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	"time"
)

type (
	ClientDto struct {
		ID    int    `json:"id"`
		Cpf   int    `json:"cpf"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	Order struct {
		ID          int              `json:"id" gorm:"primaryKey;autoIncrement"`
		ClientId    *int             `json:"client_id"`
		StatusOrder enum.OrderStatus `json:"status_order"`
		Amount      float32          `json:"amount"`
		CreatedAt   time.Time        `json:"created_at,omitempty"`
		Products    []*Product       `json:"products" gorm:"many2many:orders_products"`
	}

	OrderDto struct {
		ID          int          `json:"id"`
		ClientID    *int         `json:"client_id"`
		StatusOrder string       `json:"status_order"`
		Amount      float32      `json:"amount"`
		CreatedAt   time.Time    `json:"created_at"`
		Products    []ProductDto `json:"products"`
	}

	ProductDto struct {
		ID          int     `json:"id"`
		NameProduct string  `json:"name_product"`
		Description string  `json:"description"`
		Price       float32 `json:"price"`
		CategoryId  int     `json:"category_id"`
		Active      bool    `json:"active"`
	}

	Product struct {
		ID          int     `json:"id" gorm:"primaryKey;autoIncrement"`
		NameProduct string  `json:"name_product"`
		Description string  `json:"description"`
		Price       float32 `json:"price"`
		CategoryId  int     `json:"category_id"`
		Active      bool    `json:"active"`
	}
)
