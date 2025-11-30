package http

import (
	authapp "go-ai/internal/application/auth"
	restaurantapp "go-ai/internal/application/restaurant"
	"go-ai/internal/infra/cache"
	authrepo "go-ai/internal/infra/db/auth"
	restaurantrepo "go-ai/internal/infra/db/restaurant"
	"go-ai/internal/infra/storage"
	"go-ai/internal/transport/http/handler"
	"go-ai/internal/transport/http/middlewares"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func Router(pool *pgxpool.Pool, e *echo.Echo, redis *redis.Client) {
	api := e.Group("/api")
	// ---- API AUTH ----
	authRepo := authrepo.NewAuthRepo(pool)
	authCache := cache.NewAuthCache(redis)
	authMiddleware := middlewares.NewAuthMiddleware(authCache)
	registerUC := authapp.NewRegisterUseCase(authRepo, authCache)
	loginUC := authapp.NewLoginUseCase(authRepo, authCache)
	refreshUC := authapp.NewRefreshTokenUseCase(authRepo, authCache)
	profileUC := authapp.NewGetProfileUseCase(authRepo, authCache)
	authHandler := handler.NewAuthHandler(
		registerUC,
		loginUC,
		refreshUC,
		profileUC,
	)
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh-token", authHandler.RefreshToken)
		authGroup.GET("/profile", authHandler.GetProfile, authMiddleware.Handle)
	}

	minioClient := storage.NewMinioClient()
	uploadHander := handler.NewUploadHandler(
		minioClient,
	)
	uploadGroup := api.Group("/upload")
	{
		uploadGroup.POST("/logo", uploadHander.UploadLogoHandler())
	}

	restaurantRepo := restaurantrepo.NewRestaurantRepo(pool)
	createRestaurantUC := restaurantapp.NewCreateRestaurantUseCase(restaurantRepo)
	restaurantHandler := handler.NewRestaurantHandler(createRestaurantUC)

	restaurantGroup := api.Group("/restaurant")
	{
		restaurantGroup.POST("/", restaurantHandler.Create, authMiddleware.Handle)
	}
}
