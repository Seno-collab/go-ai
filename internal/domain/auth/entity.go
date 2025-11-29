package auth

import (
	"go-ai/pkg/common"
	"strings"

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

func (u *Entity) Validate() error {
	if !strings.Contains(u.Email, "@") {
		return common.ErrInvalidEmail
	}
	if strings.TrimSpace(u.FullName) == "" {
		return common.ErrInvalidName
	}
	if strings.TrimSpace(u.Password) == "" {
		return common.ErrInvalidPassword
	}
	return nil
}

func (u *Entity) ValidateLogin() error {
	if !strings.Contains(u.Email, "@") {
		return common.ErrInvalidEmail
	}
	if strings.TrimSpace(u.Password) == "" {
		return common.ErrInvalidPassword
	}
	return nil
}
