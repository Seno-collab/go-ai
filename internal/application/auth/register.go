package authapp

import (
	"database/sql"
	"errors"
	"go-ai/internal/domain/auth"
	"go-ai/internal/infra/cache"
	"go-ai/internal/transport/http/status"
	uilts "go-ai/pkg/utils"
	"strings"

	"github.com/google/uuid"
)

type RegisterUseCase struct {
	repo  auth.Repository
	cache *cache.AuthCache
}

func NewRegisterUseCase(repo auth.Repository, cache *cache.AuthCache) *RegisterUseCase {
	return &RegisterUseCase{
		repo:  repo,
		cache: cache,
	}
}

func (s *RegisterUseCase) Execute(request RegisterRequest) (uuid.UUID, error) {
	if !strings.Contains(request.Email, "@") {
		return uuid.Nil, status.ErrInvalidEmail
	}
	if request.FullName == "" {
		return uuid.Nil, status.ErrInvalidName
	}
	if request.Password == "" {
		return uuid.Nil, status.ErrInvalidPassword
	}
	// unique email
	record, err := s.repo.GetByEmail(request.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, err
		}
	}
	if record != nil {
		return uuid.Nil, status.ErrUserAlreadyExists
	}
	record, err = s.repo.GetByName(request.FullName)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, err
		}
	}
	if record != nil {
		return uuid.Nil, status.ErrNameAlreadyExists
	}
	hasedPassword, err := uilts.HashPassword(request.Password)
	if err != nil {
		return uuid.Nil, status.ErrInternalServerError
	}
	return s.repo.CreateUser(&auth.Entity{
		FullName: request.FullName,
		Email:    request.Email,
		Password: hasedPassword,
	})
}
