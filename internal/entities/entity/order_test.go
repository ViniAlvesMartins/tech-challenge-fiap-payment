package entity

import (
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOrder_SetAmount(t *testing.T) {
	t.Run("test set amount successfully", func(t *testing.T) {
		o := Order{
			ID:          1,
			ClientId:    nil,
			StatusOrder: enum.AWAITING_PAYMENT,
			Amount:      123.45,
			CreatedAt:   time.Now(),
		}

		o.SetAmount(678.9)

		assert.Equal(t, o.Amount, float32(678.9))
	})
}
