package response_payment_service

type CreateQRCode struct {
	QrData string `json:"qr_data"`
}
