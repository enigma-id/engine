package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/enigma-id/engine"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var (
	Server     *echo.Echo // Global Echo server instance
	RestConfig *Config    // Configuration for REST server
	Logger     *zap.Logger
)

// Config holds configuration settings for the REST server.
type Config struct {
	Server    string
	IsDev     bool
	JwtSecret string
}

// NewServer initializes a new Echo server with middleware and routes.
func NewServer(config *Config, registerRoutes func(e *echo.Echo)) {
	Server = echo.New()
	RestConfig = config

	// Basic server configuration
	Server.Debug = config.IsDev
	Server.HidePort = true
	Server.HideBanner = true
	// Server.Binder = &CustomBinder{}
	// Server.HTTPErrorHandler = CustomHTTPErrorHandler

	// Middleware stack
	Server.Use(HTTPLogger())
	Server.Use(middleware.Recover())
	Server.Use(middleware.CORS())
	Server.Use(middleware.Gzip())
	Server.Use(middleware.RequestID())
	Server.Use(middleware.AddTrailingSlash())
	Server.Use(middleware.BodyDump(RequestFailureDumper))

	// Custom route registration
	registerRoutes(Server)

	// Health check endpoint
	Server.GET("/health", handlerHealth)

	// Dev-only endpoints (e.g. Swagger docs)
	if RestConfig.IsDev {
		Server.File("/swagger.json", "docs/swagger.json")
	}

	logRoutes(Server)
}

// Start runs the Echo server.
func Start() error {
	Logger.Info(fmt.Sprintf("Starting Rest Server: %s", RestConfig.Server))
	return Server.Start(RestConfig.Server)
}

// Shutdown gracefully shuts down the server.
func Shutdown(err error) {
	Server.Shutdown(context.Background())
}

// logRoutes prints registered routes in development mode.
func logRoutes(e *echo.Echo) {
	if !Server.Debug {
		return
	}

	fmt.Println("\nREGISTERED ROUTES: ")
	fmt.Println(strings.Repeat("-", 120))
	fmt.Printf("%-10s | %-50s | %-54s\n", "METHOD", "URL PATH", "REQ. HANDLER")
	fmt.Println(strings.Repeat("-", 120))

	for _, r := range e.Routes() {
		if !strings.HasSuffix(r.Path, "*") {
			fmt.Printf("%-10s | %-50s | %-54s\n", r.Method, r.Path, r.Name)
		}
	}

	fmt.Println(strings.Repeat("-", 120))
}

// handlerHealth responds with basic health check info.
func handlerHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"service": engine.ServiceConfig.Name,
		"version": engine.ServiceConfig.Version,
		"host":    os.Getenv("DAP_AGENT"),
		"time":    time.Now(),
	})
}

// IDFromContext extracts an integer ID from the Echo context params.
func IDFromContext(c echo.Context) int64 {
	id, _ := strconv.Atoi(c.Param("id"))
	return int64(id)
}
