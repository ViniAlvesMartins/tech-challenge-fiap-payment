package entity

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExternalPayment(t *testing.T) {
	u := uuid.New()

	e := ExternalPayment{
		ID:      &u,
		OrderID: 1,
		Type:    enum.QRCODE,
		Amount:  123.45,
	}

	assert.IsType(t, e, ExternalPayment{})
}
