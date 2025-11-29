package auth_service

import (
	"database/sql"
	"errors"
	"fmt"
	"go-ai/internal/config"
	"go-ai/internal/domain/auth"
	"go-ai/internal/infra/cache"
	"go-ai/pkg/common"
	uilts "go-ai/pkg/utils"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo  auth.Repository
	cache *cache.AuthCache
}

func NewAuthService(repo auth.Repository, cache *cache.AuthCache) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

func (s *Service) Register(u *auth.Entity) (uuid.UUID, error) {
	if err := u.Validate(); err != nil {
		return uuid.Nil, err
	}

	// unique email
	record, err := s.repo.GetByEmail(u.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, err
		}
	}
	if record != nil {
		return uuid.Nil, common.ErrUserAlreadyExists
	}
	record, err = s.repo.GetByName(u.FullName)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, err
		}
	}
	if record != nil {
		return uuid.Nil, common.ErrNameAlreadyExists
	}
	hasedPassword, err := uilts.HashPassword(u.Password)
	if err != nil {
		return uuid.Nil, common.ErrInternalServerError
	}
	u.Password = hasedPassword
	return s.repo.CreateUser(u)
}

func (s *Service) Login(u *auth.Entity) (string, string, int, error) {
	config, _ := config.LoadConfig()

	if err := u.ValidateLogin(); err != nil {
		return "", "", 0, err
	}
	storedUser, err := s.repo.GetByEmail(u.Email)
	if err != nil {
		return "", "", 0, err
	}
	if storedUser.IsActive == false {
		return "", "", 0, common.ErrUserInactive
	}
	if !uilts.CheckPasswordHash(u.Password, storedUser.Password) {
		return "", "", 0, auth.ErrPasswordVerifyFail
	}
	accessToken, err := uilts.GenerateToken(storedUser.ID, storedUser.Email, storedUser.Role, config.JwtAccessSecret, config.JwtExpiresIn)
	if err != nil {
		return "", "", 0, auth.ErrTokenGenerateFail
	}
	refreshToken, err := uilts.GenerateToken(storedUser.ID, storedUser.Email, storedUser.Role, config.JwtRefreshSecret, config.JwtRefreshExpiresIn)
	if err != nil {
		return "", "", 0, auth.ErrTokenGenerateFail
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
	return accessToken, refreshToken, config.JwtExpiresIn, nil
}

func (s *Service) RefreshToken(refreshToken string) (string, string, int, error) {
	if refreshToken == "" {
		return "", "", 0, auth.ErrorRefreshTokenEmpty
	}
	config, _ := config.LoadConfig()
	claims, err := uilts.VerifyToken(refreshToken, config.JwtRefreshSecret)
	if err != nil {
		return "", "", 0, auth.ErrTokenNotActive
	}
	userId := claims.UserId
	if userId == uuid.Nil {
		return "", "", 0, auth.ErrTokenMissing
	}
	email := claims.Email
	if email == "" {
		return "", "", 0, auth.ErrTokenMissing
	}
	role := claims.Role
	if role == "" {
		return "", "", 0, auth.ErrTokenMissing
	}
	keyRefreshToken := fmt.Sprintf("refresh_token_%s", userId)
	cachedRefreshToken, err := s.cache.GetRefreshTokenCache(keyRefreshToken)
	if err != nil {
		return "", "", 0, err
	}
	if cachedRefreshToken != refreshToken {
		return "", "", 0, auth.ErrTokenMalformed
	}
	accessToken, err := uilts.GenerateToken(userId, email, role, config.JwtAccessSecret, config.JwtExpiresIn)
	if err != nil {
		return "", "", 0, auth.ErrTokenGenerateFail
	}
	newRefreshToken, err := uilts.GenerateToken(userId, email, role, config.JwtRefreshSecret, config.JwtRefreshExpiresIn)
	if err != nil {
		return "", "", 0, auth.ErrTokenGenerateFail
	}
	record, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", "", 0, err
	}
	dataCache := &cache.AuthData{
		UserId:   record.ID,
		Role:     record.Role,
		Email:    record.Email,
		IsActive: record.IsActive,
		FullName: record.FullName,
	}
	keyAuthCache := fmt.Sprintf("profile_%s", record.ID.String())
	s.cache.SetAuthCache(keyAuthCache, dataCache, time.Duration(config.JwtExpiresIn*int(time.Second)))
	s.cache.SetRefreshTokenCache(keyRefreshToken, newRefreshToken, time.Duration(config.JwtRefreshExpiresIn*int(time.Second)))
	return accessToken, newRefreshToken, config.JwtExpiresIn, nil
}

type Profile struct {
	Email    string
	FullName string
	Role     string
	IsActive bool
}

func (s *Service) GetProfile(userId uuid.UUID) (*Profile, error) {
	keyAuth := fmt.Sprintf("profile_%s", userId.String())
	cacheData, err := s.cache.GetAuthCache(keyAuth)
	if err != nil {
		return nil, err
	}
	profile := &Profile{}
	if cacheData == nil {
		record, err := s.repo.GetById(userId)
		if err != nil {
			return nil, err
		}
		profile = &Profile{
			Email:    record.Email,
			FullName: record.FullName,
			Role:     record.Role,
			IsActive: record.IsActive,
		}
		authData := &cache.AuthData{
			Email:    record.Email,
			FullName: record.FullName,
			Role:     record.Role,
			IsActive: record.IsActive,
		}
		s.cache.SetAuthCache(keyAuth, authData, time.Duration(60*int(time.Minute)))
		return profile, nil
	}
	profile = &Profile{
		Email:    cacheData.Email,
		FullName: cacheData.FullName,
		Role:     cacheData.Role,
		IsActive: cacheData.IsActive,
	}
	return profile, nil
}
