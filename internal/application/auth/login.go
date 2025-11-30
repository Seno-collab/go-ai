package authapp

import (
	"fmt"
	"go-ai/internal/config"
	"go-ai/internal/domain/auth"
	"go-ai/internal/infra/cache"
	"go-ai/internal/transport/http/status"

	uilts "go-ai/pkg/utils"
	"time"
)

type LoginUseCase struct {
	repo  auth.Repository
	cache *cache.AuthCache
}

func NewLoginUseCase(repo auth.Repository, cache *cache.AuthCache) *LoginUseCase {
	return &LoginUseCase{
		repo:  repo,
		cache: cache,
	}
}

func (s *LoginUseCase) Execute(request LoginRequest) (*LoginResponse, error) {
	config, _ := config.LoadConfig()

	if request.Email == "" {
		return nil, status.ErrInvalidEmail
	}

	if request.Password == "" {
		return nil, status.ErrInvalidPassword
	}

	storedUser, err := s.repo.GetByEmail(request.Email)
	if err != nil {
		return nil, err
	}
	if storedUser.IsActive == false {
		return nil, status.ErrUserInactive
	}
	if !uilts.CheckPasswordHash(request.Password, storedUser.Password) {
		return nil, auth.ErrPasswordVerifyFail
	}
	accessToken, err := uilts.GenerateToken(storedUser.ID, storedUser.Email, storedUser.Role, config.JwtAccessSecret, config.JwtExpiresIn)
	if err != nil {
		return nil, auth.ErrTokenGenerateFail
	}
	refreshToken, err := uilts.GenerateToken(storedUser.ID, storedUser.Email, storedUser.Role, config.JwtRefreshSecret, config.JwtRefreshExpiresIn)
	if err != nil {
		return nil, auth.ErrTokenGenerateFail
	}
	dataCache := &cache.AuthData{
		UserId:   storedUser.ID,
		Role:     storedUser.Role,
		Email:    storedUser.Email,
		IsActive: storedUser.IsActive,
		FullName: storedUser.FullName,
	}
	keyAuthCache := fmt.Sprintf("profile_%s", storedUser.ID.String())
	s.cache.SetAuthCache(keyAuthCache, dataCache, time.Duration(config.JwtExpiresIn*int(time.Second)))
	keyRefreshToken := fmt.Sprintf("refresh_token_%s", storedUser.ID.String())
	s.cache.SetRefreshTokenCache(keyRefreshToken, refreshToken, time.Duration(config.JwtRefreshExpiresIn*int(time.Second)))
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpresIn:     config.JwtExpiresIn,
	}, nil
}
