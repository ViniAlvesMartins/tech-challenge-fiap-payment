package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient(t *testing.T) {
	c := Client{
		ID:    1,
		Cpf:   123456789,
		Name:  "Person McPerson",
		Email: "mcperson@example.com",
	}

	assert.IsType(t, c, Client{})
}
