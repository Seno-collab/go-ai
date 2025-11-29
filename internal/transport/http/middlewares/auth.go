package middlewares

import (
	"fmt"
	"go-ai/internal/config"
	"go-ai/internal/infra/cache"
	"strings"
	"time"

	"go-ai/pkg/common"
	uilts "go-ai/pkg/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	Cache *cache.AuthCache
}

func NewAuthMiddleware(cache *cache.AuthCache) *AuthMiddleware {
	return &AuthMiddleware{
		Cache: cache,
	}
}

func (m *AuthMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		config, _ := config.LoadConfig()
		// Implement authentication logic here
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return common.ErrorResponse(c, 401, "Missing Authorization header")
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return common.ErrorResponse(c, 401, "Invalid Authorization header format")
		}
		token := parts[1]
		claims, err := uilts.VerifyToken(token, config.JwtAccessSecret)
		if err != nil || claims == nil {
			return common.ErrorResponse(c, 401, "Invalid token")
		}
		exp := claims.ExpiresAt
		if exp == nil {
			return common.ErrorResponse(c, 401, "Token has expired")
		}
		if time.Now().After(exp.Time) {
			return common.ErrorResponse(c, 401, "Token has expired")
		}
		userId := claims.UserId
		if userId == uuid.Nil {
			return common.ErrorResponse(c, 401, "Unauthorized access")
		}
		keyAuth := fmt.Sprintf("profile_%s", claims.UserId.String())
		authData, err := m.Cache.GetAuthCache(keyAuth)
		if err != nil || authData == nil {
			return common.ErrorResponse(c, 401, "Unauthorized access")
		}
		c.Set("user_id", claims.UserId)
		return next(c)
	}
}
