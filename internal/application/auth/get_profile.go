package authapp

import (
	"fmt"
	"go-ai/internal/domain/auth"
	"go-ai/internal/infra/cache"
	"time"

	"github.com/google/uuid"
)

type GetProfileUseCase struct {
	repo  auth.Repository
	cache *cache.AuthCache
}

func NewGetProfileUseCase(repo auth.Repository, cache *cache.AuthCache) *GetProfileUseCase {
	return &GetProfileUseCase{
		repo:  repo,
		cache: cache,
	}
}

func (uc *GetProfileUseCase) Execute(userId uuid.UUID) (*GetProfileResponse, error) {
	keyAuth := fmt.Sprintf("profile_%s", userId.String())
	cacheData, err := uc.cache.GetAuthCache(keyAuth)
	if err != nil {
		return nil, err
	}
	profile := &GetProfileResponse{}
	if cacheData == nil {
		record, err := uc.repo.GetById(userId)
		if err != nil {
			return nil, err
		}
		profile = &GetProfileResponse{
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
		uc.cache.SetAuthCache(keyAuth, authData, time.Duration(60*int(time.Minute)))
		return profile, nil
	}
	profile = &GetProfileResponse{
		Email:    cacheData.Email,
		FullName: cacheData.FullName,
		Role:     cacheData.Role,
		IsActive: cacheData.IsActive,
	}
	return profile, nil
}
