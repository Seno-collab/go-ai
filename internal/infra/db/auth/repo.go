package authrepo

import (
	"context"
	auth "go-ai/internal/domain/auth"
	sqlc "go-ai/internal/infra/sqlc/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	q *sqlc.Queries
}

func NewAuthRepo(pool *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		q: sqlc.New(pool),
	}
}

func (au *AuthRepo) GetByEmail(email string) (*auth.Entity, error) {
	u, err := au.q.GetUserByEmail(context.Background(), &email)
	if err != nil {
		return nil, err
	}

	return &auth.Entity{
		ID:       u.ID,
		Email:    *u.Email,
		FullName: u.FullName,
		Password: u.PasswordHash,
		Role:     *u.RoleName,
		IsActive: u.IsActive,
	}, nil
}

func (au *AuthRepo) CreateUser(a *auth.Entity) (uuid.UUID, error) {
	id, err := au.q.CreateUser(context.Background(), sqlc.CreateUserParams{
		Email:        &a.Email,
		PasswordHash: a.Password,
		FullName:     a.FullName,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
func (au *AuthRepo) GetByName(name string) (*auth.Entity, error) {
	u, err := au.q.GetUserByName(context.Background(), name)
	if err != nil {
		return nil, err
	}
	return &auth.Entity{
		ID:       u.ID,
		Email:    *u.Email,
		FullName: u.FullName,
		Role:     *u.RoleName,
		IsActive: u.IsActive,
	}, nil
}

func (au *AuthRepo) GetById(id uuid.UUID) (*auth.Entity, error) {
	u, err := au.q.GetUserByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return &auth.Entity{
		ID:       u.ID,
		Email:    *u.Email,
		FullName: u.FullName,
		Role:     *u.RoleName,
		IsActive: u.IsActive,
	}, nil
}
