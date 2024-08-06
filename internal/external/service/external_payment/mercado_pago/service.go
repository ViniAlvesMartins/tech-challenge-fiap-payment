package mercado_pago

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
)

type PaymentGateway struct {
}

type Response struct {
	QRData         string `json:"qr_data"`
	InStoreOrderID string `json:"in_store_order_id"`
}

func NewPaymentGateway() *PaymentGateway {
	return &PaymentGateway{}
}

func (g *PaymentGateway) CreateQRCode(p entity.Payment) (Response, error) {
	return Response{
		QRData:         "00020101021243650016COM.MERCADOLIBRE02013063638f1192a-5fd1-4180-a180-8bcae3556bc35204000053039865802BR5925IZABEL AAAA DE MELO6007BARUERI62070503***63040B6D",
		InStoreOrderID: "1",
	}, nil
}
