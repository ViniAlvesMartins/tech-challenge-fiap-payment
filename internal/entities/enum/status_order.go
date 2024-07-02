package enum

import (
	"slices"
)

type StatusOrder string

const (
	AWAITING_PAYMENT StatusOrder = "AWAITING_PAYMENT"
	RECEIVED         StatusOrder = "RECEIVED"
	PREPARING        StatusOrder = "PREPARING"
	READY            StatusOrder = "READY"
	FINISHED         StatusOrder = "FINISHED"
)

func ValidateStatus(val string) bool {
	validStatus := []StatusOrder{AWAITING_PAYMENT, RECEIVED, PREPARING, READY, FINISHED}
	return slices.Contains(validStatus, StatusOrder(val))
}
