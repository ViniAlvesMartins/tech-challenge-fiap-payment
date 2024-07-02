//go:generate mockgen -destination=mock/uuid.go -source=uuid.go -package=mock
package uuid

import "github.com/google/uuid"

type Interface interface {
	New() uuid.UUID
	NewString() string
}

type UUID struct{}

func (u *UUID) New() uuid.UUID {
	return uuid.New()
}

func (u *UUID) NewString() string {
	return u.New().String()
}
