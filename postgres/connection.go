package postgres

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/zap"
)

// Config holds the configuration for the Postgres connection.
type Config struct {
	Server     string
	Username   string
	Password   string
	Database   string
	Datasource string
	Debug      bool
}

var db *bun.DB // internal DB instance

// setDefault builds the DSN if not explicitly provided.
func (c *Config) setDefault() {
	if c.Datasource == "" {
		c.Datasource = fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			c.Username, c.Password, c.Server, c.Database,
		)
	}
}

// NewConnection initializes a new Postgres connection and sets the global DB.
func NewConnection(c *Config, logger *zap.Logger) error {
	c.setDefault()

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(c.Datasource)))

	if err := sqldb.Ping(); err != nil {
		return err
	}

	db = bun.NewDB(sqldb, pgdialect.New())

	if c.Debug {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
		))
	}

	logger.Info(fmt.Sprintf("Connected to Postgres Server: %s@%s", c.Server, c.Database))
	return nil
}

// GetDB returns the bun.DB instance.
func GetDB() *bun.DB {
	return db
}
