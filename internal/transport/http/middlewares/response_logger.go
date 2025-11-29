package middlewares

import (
	"bytes"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

const maxLoggedBodyBytes = 2048

func ResponseLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			res := c.Response()
			req := c.Request()

			logger := zerolog.Ctx(req.Context())

			latency := time.Since(start)

			var responseBuf bytes.Buffer
			origWriter := res.Writer
			respLogger := &responseCaptureWriter{
				ResponseWriter: origWriter,
				buf:            &responseBuf,
			}
			res.Writer = respLogger
			err := next(c)
			res.Writer = origWriter

			responseBody := formatBody(responseBuf.Bytes())
			logger.Info().
				Int("status", res.Status).
				Int64("latency_ms", latency.Milliseconds()).
				Int64("size", res.Size).
				Str("response_body", responseBody).Msg("response sent")
			return err
		}
	}

}

func formatBody(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	if len(data) > maxLoggedBodyBytes {
		data = data[:maxLoggedBodyBytes]
	}

	body := strings.TrimSpace(string(data))
	if body == "" {
		return ""
	}
	return body
}

type responseCaptureWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (w *responseCaptureWriter) Write(b []byte) (int, error) {
	if len(b) > 0 {
		if w.buf.Len()+len(b) <= maxLoggedBodyBytes {
			w.buf.Write(b)
		} else {
			remaining := maxLoggedBodyBytes - w.buf.Len()
			if remaining > 0 {
				w.buf.Write(b[:remaining])
			}
		}
	}
	return w.ResponseWriter.Write(b)
}

func (w *responseCaptureWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
