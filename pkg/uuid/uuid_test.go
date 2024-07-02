package uuid

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUUID_New(t *testing.T) {
	u := UUID{}

	assert.IsType(t, uuid.UUID{}, u.New())
}

func TestUUID_NewString(t *testing.T) {
	u := UUID{}

	assert.NotEmpty(t, u.NewString())
}
