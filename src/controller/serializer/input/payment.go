package input

type PaymentDto struct {
	Type    string `json:"type" validate:"required" error:"Tipo de pagamento é obrigatorio"`
	OrderId int    `json:"orderId" validate:"required" error:"order id é obrigatorio"`
}
