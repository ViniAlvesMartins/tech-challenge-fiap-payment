package enum

type PaymentType string

type PaymentStatus string

const (
	CREDIT PaymentType = "CREDIT"
	DEBIT  PaymentType = "DEBIT"
	CASH   PaymentType = "CASH"
	PIX    PaymentType = "PIX"
	QRCODE PaymentType = "QRCODE"

	PENDING   PaymentStatus = "PENDING"
	CONFIRMED PaymentStatus = "CONFIRMED"
	CANCELED  PaymentStatus = "CANCELED"
)
