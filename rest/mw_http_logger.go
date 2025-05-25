package rest

import (
	"fmt"
	"time"

	"github.com/enigma-id/engine"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// HTTPLogger returns a middleware that logs HTTP requests.
func HTTPLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			return logRequest(next, c)
		}
	}
}

// logRequest print all http request on consoles.
func logRequest(hand echo.HandlerFunc, c echo.Context) (err error) {
	start := time.Now()
	req := c.Request()
	res := c.Response()
	if err = hand(c); err != nil {
		c.Error(err)
	}
	end := time.Now()
	latency := end.Sub(start) / 1e5

	id := req.Header.Get(echo.HeaderXRequestID)
	if id == "" {
		id = res.Header().Get(echo.HeaderXRequestID)
	}

	var fields = []zap.Field{
		zap.String("path", req.URL.Path),
		zap.String("id", id),
		zap.String("query", req.URL.RawQuery),
		zap.String("ip", c.RealIP()),
		zap.String("user-agent", req.UserAgent()),
		zap.String("latecy", fmt.Sprintf("%1.1fms", float64(latency))),
	}

	if err == nil {
		engine.Logger.Info(fmt.Sprintf("%s/%d", req.Method, res.Status), fields...)
	} else {
		engine.Logger.Warn(fmt.Sprintf("%s/%d", req.Method, res.Status), fields...)
	}

	return
}
