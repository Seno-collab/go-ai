package auth

import (
	"github.com/google/uuid"
)

type Entity struct {
	ID       uuid.UUID
	FullName string
	Email    string
	Password string
	Role     string
	IsActive bool
}
