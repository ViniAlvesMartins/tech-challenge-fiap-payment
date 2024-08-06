package enum

import "slices"

type (
	PaymentType   string
	PaymentStatus string
	OrderStatus   string
)

const (
	PaymentTypeQRCode PaymentType = "QRCODE"

	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusConfirmed PaymentStatus = "CONFIRMED"
	PaymentStatusCanceled  PaymentStatus = "CANCELED"

	OrderStatusAwaitingPayment OrderStatus = "AWAITING_PAYMENT"
	OrderStatusReceived        OrderStatus = "RECEIVED"
	OrderStatusPreparing       OrderStatus = "PREPARING"
	OrderStatusReady           OrderStatus = "READY"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusFinished        OrderStatus = "FINISHED"
)

func ValidateStatus(val string) bool {
	validStatus := []OrderStatus{OrderStatusAwaitingPayment, OrderStatusReceived, OrderStatusPreparing, OrderStatusReady, OrderStatusFinished}
	return slices.Contains(validStatus, OrderStatus(val))
}
