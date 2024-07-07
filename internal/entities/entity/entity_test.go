package entity

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

func TestClient(t *testing.T) {
	c := Client{
		ID:    1,
		Cpf:   123456789,
		Name:  "Person McPerson",
		Email: "mcperson@example.com",
	}

	assert.IsType(t, c, Client{})
}

func TestOrder(t *testing.T) {
	t.Run("test set amount successfully", func(t *testing.T) {
		o := Order{
			ID:          1,
			ClientId:    nil,
			StatusOrder: enum.AWAITING_PAYMENT,
			Amount:      123.45,
			CreatedAt:   time.Now(),
		}

		assert.Equal(t, o.Amount, float32(123.45))
	})
}
