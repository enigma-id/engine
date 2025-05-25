package rest

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RequestFailureDumper logs failed HTTP requests and their payloads for easier debugging.
// It is triggered only when the response status is not 200 OK.
func RequestFailureDumper(c echo.Context, reqBody, resBody []byte) {
	// Skip logging if response was successful
	if c.Response().Status == http.StatusOK {
		return
	}

	req := c.Request()
	res := c.Response()

	// Extract request ID from request or response headers
	requestID := req.Header.Get(echo.HeaderXRequestID)
	if requestID == "" {
		requestID = res.Header().Get(echo.HeaderXRequestID)
	}

	// Parse request and response bodies as JSON (best effort)
	var requestPayload = json.RawMessage(reqBody)
	var responsePayload = json.RawMessage(resBody)

	// Collect structured log fields
	fields := []zap.Field{
		zap.String("method", req.Method),
		zap.String("path", req.URL.Path),
		zap.Int("status", res.Status),
		zap.String("request_id", requestID),
		zap.String("query", req.URL.RawQuery),
		zap.Any("request", &requestPayload),
		zap.Any("response", &responsePayload),
	}

	// Log with structured context
	Logger.Warn("HTTP request failed", fields...)
}
