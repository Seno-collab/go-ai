package middlewares

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func RequestLoggerMiddleware(baseLogger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			requestID := res.Header().Get(echo.HeaderXRequestID)
			var requestBody []byte
			requestBodyLogged := false
			if req.Body != nil {
				bodyBytes, err := io.ReadAll(req.Body)
				if err == nil {
					if len(bodyBytes) > 0 {
						requestBodyLogged = true
					}
					req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				} else {
					requestBodyLogged = true
				}
				requestBody = bodyBytes
			}
			logger := baseLogger.With().
				Str("request_id", requestID).
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Str("remote_ip", c.RealIP()).
				Str("user_agent", req.UserAgent())
			if requestBodyLogged {
				logger = logger.Str("request_body", prettyJSON(requestBody))
			}
			reqLogger := logger.Logger()
			reqLogger.Info().Msg("incoming request")

			ctx := baseLogger.With().Str("request_id", requestID).Logger().WithContext(req.Context())
			c.SetRequest(req.WithContext(ctx))

			return next(c)
		}
	}
}

func prettyJSON(raw []byte) string {
	var out bytes.Buffer
	err := json.Compact(&out, raw)
	if err != nil {
		return string(raw)
	}
	return out.String()
}
