package middlewares

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func RequestIDMiddleware(baseLogger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			req := c.Request()
			res := c.Response()

			requestID := req.Header.Get(echo.HeaderXRequestID)
			res.Header().Set(echo.HeaderXRequestID, requestID)
			if requestID == "" {
				requestID = uuid.New().String()
				c.Response().Header().Set(echo.HeaderXRequestID, requestID)
			}
			reqLogger := baseLogger.With().
				Str("request_id", requestID).
				Logger()

			ctx := reqLogger.WithContext(req.Context())
			c.SetRequest(req.WithContext(ctx))

			return next(c)
		}
	}
}
