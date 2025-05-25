package engine

import (
	"github.com/enigma-id/engine/log"
	"github.com/oklog/run"
	"go.uber.org/zap"
)

type Config struct {
	Name    string
	Version string
	Host    string
	IsDev   bool
}

var (
	ServiceConfig *Config
	Logger        *zap.Logger
	Routine       run.Group
)

// Start initializes the service configuration and logger.
func Start(cfg *Config) *Config {
	ServiceConfig = cfg
	Logger = NewLogger(cfg.Name)
	return ServiceConfig
}

// NewLogger creates a named logger using the global config.
func NewLogger(name string) *zap.Logger {
	return log.BuildLogger(ServiceConfig.IsDev).Named(name)
}
