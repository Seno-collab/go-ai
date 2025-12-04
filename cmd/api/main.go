package main

import (
	"context"
	"fmt"
	_ "go-ai/docs"
	"go-ai/internal/config"
	"go-ai/internal/infra/cache"
	"go-ai/internal/infra/db"
	httpHandler "go-ai/internal/transport/http"
	"go-ai/internal/transport/http/middlewares"
	"go-ai/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/time/rate"
)

// go ai
// mission using golang build ai. Development application AI

func main() {
	e := echo.New()
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	logger := logger.NewLogger()
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))
	e.Use(middlewares.RequestIDMiddleware(logger))
	e.Use(middleware.ContextTimeout(60 * time.Second))
	e.Use(middleware.Recover())
	e.Use(middlewares.RequestLoggerMiddleware(logger))
	e.Use(middlewares.ResponseLoggerMiddleware())
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("Load config error")
		return
	}
	logger.Info().
		Str("environment", cfg.Environment).
		Str("server_port", cfg.ServerPort).
		Str("db_host", cfg.DBHost).
		Msg("configuration loaded successfully")
	dsnPg := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	pool, err := db.ConnectPostgres(dsnPg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Connect db fail")
		return
	}

	dsnRedis := fmt.Sprintf("redis://%s:%s@%s:%d/%d", "", cfg.RedisPassword, cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
	redisClient, err := cache.ConnectRedis(dsnRedis)
	if err != nil {
		logger.Fatal().Fields(map[string]any{
			"Error": err,
		}).Msg("Connect redis fail")
	}

	port := fmt.Sprintf(":%s", cfg.ServerPort)
	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Shutting down the server")
		}
	}()
	httpHandler.Router(pool, e, redisClient)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}
	logger.Info().Msg("Server exited gracefully")
}
