package authapp

import (
	"fmt"
	"go-ai/internal/config"
	"go-ai/internal/domain/auth"
	"go-ai/internal/infra/cache"
	"go-ai/pkg/utils"
	uilts "go-ai/pkg/utils"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenUseCase struct {
	repo  auth.Repository
	cache *cache.AuthCache
}

func NewRefreshTokenUseCase(repo auth.Repository, cache *cache.AuthCache) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		repo:  repo,
		cache: cache,
	}
}

func (uc *RefreshTokenUseCase) Execute(request RefreshTokenRequest) (*RefreshTokenResponse, error) {
	if request.RefreshToken == "" {
		return nil, auth.ErrTokenInvalid
	}
	config, _ := config.LoadConfig()
	claims, err := utils.VerifyToken(request.RefreshToken, config.JwtRefreshSecret)
	if err != nil {
		return nil, auth.ErrTokenNotActive
	}
	userId := claims.UserId
	if userId == uuid.Nil {
		return nil, auth.ErrTokenMissing
	}
	email := claims.Email
	if email == "" {
		return nil, auth.ErrTokenMissing
	}
	role := claims.Role
	if role == "" {
		return nil, auth.ErrTokenMissing
	}
	keyRefreshToken := fmt.Sprintf("refresh_token_%s", userId)
	cachedRefreshToken, err := uc.cache.GetRefreshTokenCache(keyRefreshToken)
	if err != nil {
		return nil, err
	}
	if cachedRefreshToken != request.RefreshToken {
		return nil, auth.ErrTokenMalformed
	}
	accessToken, err := uilts.GenerateToken(userId, email, role, config.JwtAccessSecret, config.JwtExpiresIn)
	if err != nil {
		return nil, auth.ErrTokenGenerateFail
	}
	refreshToken, err := uilts.GenerateToken(userId, email, role, config.JwtRefreshSecret, config.JwtRefreshExpiresIn)
	if err != nil {
		return nil, auth.ErrTokenGenerateFail
	}
	record, err := uc.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	dataCache := &cache.AuthData{
		UserId:   record.ID,
		Role:     record.Role,
		Email:    record.Email,
		IsActive: record.IsActive,
		FullName: record.FullName,
	}
	keyAuthCache := fmt.Sprintf("profile_%s", record.ID.String())
	uc.cache.SetAuthCache(keyAuthCache, dataCache, time.Duration(config.JwtExpiresIn*int(time.Second)))
	uc.cache.SetRefreshTokenCache(keyRefreshToken, refreshToken, time.Duration(config.JwtRefreshExpiresIn*int(time.Second)))
	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    config.JwtExpiresIn,
	}, nil

}
