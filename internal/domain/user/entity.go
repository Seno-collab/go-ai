package user

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	FullName string
	Email    string
	Password string
	RoleId   int
}

// func (u *User) Validate() error {
// 	if !strings.Contains(u.Email, "@") {
// 		return ErrInvalidEmail
// 	}
// 	if strings.TrimSpace(u.FullName) == "" {
// 		return ErrInvalidName
// 	}
// 	return nil
// }
