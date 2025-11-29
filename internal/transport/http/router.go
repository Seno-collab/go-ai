package http

import (
	"go-ai/internal/infra/cache"
	"go-ai/internal/infra/db"
	authservice "go-ai/internal/service/auth"
	"go-ai/internal/transport/http/handler"
	"go-ai/internal/transport/http/middlewares"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func Router(pool *pgxpool.Pool, e *echo.Echo, redis *redis.Client) {
	api := e.Group("/api")
	// ---- API AUTH ----
	authRepo := db.NewAuthRepo(pool)
	authCache := cache.NewAuthCache(redis)
	authMiddleware := middlewares.NewAuthMiddleware(authCache)
	authService := authservice.NewAuthService(authRepo, authCache)
	authHandler := handler.NewAuthHandler(authService)
	authGroup := api.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/refresh-token", authHandler.RefreshToken)
	authGroup.GET("/profile", authHandler.GetProfile, authMiddleware.Handle)
}
