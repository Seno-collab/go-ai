package auth

import "github.com/google/uuid"

type Repository interface {
	GetById(id uuid.UUID) (*Entity, error)
	GetByEmail(email string) (*Entity, error)
	CreateUser(u *Entity) (uuid.UUID, error)
	GetByName(name string) (*Entity, error)
}
