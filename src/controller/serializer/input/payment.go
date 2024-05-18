package input

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/enum"
)

type PaymentDto struct {
	Type    string `json:"type" validate:"required" error:"Tipo de pagamento é obrigatorio"`
	OrderId int    `json:"orderId" validate:"required" error:"order id é obrigatorio"`
}

func (p *PaymentDto) ConvertToEntity() entity.Payment {
	return entity.Payment{
		OrderID: p.OrderId,
		Type:    enum.PaymentType(p.Type),
	}
}
