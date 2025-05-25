// connection.go
package redis

import (
	"fmt"

	"github.com/enigma-id/engine"
	"github.com/gomodule/redigo/redis"
)

// Cache is the global Redis instance.
var Cache *Redis

// Config contains Redis connection settings.
type Config struct {
	Name     string // Used as key prefix
	Server   string // Redis server address
	Password string // Redis password
}

// NewConnection initializes a new Redis connection and sets the global Cache.
func NewConnection(cfg *Config) error {
	pool := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.Server, redis.DialPassword(cfg.Password))
		},
	}

	Cache = &Redis{
		Prefix: cfg.Name,
		Pool:   pool,
	}

	if err := Cache.ping(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	engine.Logger.Info(fmt.Sprintf("Connected to Redis Server: %s@%s", cfg.Server, cfg.Name))
	return nil
}
