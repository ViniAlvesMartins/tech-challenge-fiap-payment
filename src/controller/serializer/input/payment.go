package input

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/src/entities/enum"
)

type PaymentDto struct {
	Type string `json:"type" validate:"required" error:"Tipo de pagamento é obrigatorio"`
}

func (p *PaymentDto) ConvertToEntity() entity.Payment {
	return entity.Payment{
		Type: enum.PaymentType(p.Type),
	}
}
