package entity

import (
	"encoding/json"
	"testing"

	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	"github.com/stretchr/testify/assert"
)

func TestPayment_GetJSONValue(t *testing.T) {
	t.Run("get json value successfully", func(t *testing.T) {
		p := Payment{
			PaymentID:    "1",
			OrderID:      "1",
			Type:         enum.QRCODE,
			CurrentState: enum.PENDING,
			Amount:       123.45,
		}

		expectedJSON, _ := json.Marshal(p)
		JSONValue, err := p.GetJSONValue()

		assert.Equal(t, string(expectedJSON), JSONValue)
		assert.Nil(t, err)
	})
}
