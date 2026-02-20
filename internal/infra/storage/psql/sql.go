package psql

import (
	"context"
	"fmt"
	"time"

	"github.com/OddKuru/core-accounts/internal/infra/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func Connect(ctx context.Context, cfg config.Psql) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Database, cfg.Username, cfg.Password, cfg.SSLMode,
	)

	conf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrapf(err, "config parse failed: %s", connString)
	}

	if cfg.MaxConnections > 0 {
		conf.MaxConns = cfg.MaxConnections
	}
	if cfg.MinConnections > 0 {
		conf.MinConns = cfg.MinConnections
	}
	if cfg.MaxConnLifetime > 0 {
		conf.MaxConnLifetime = cfg.MaxConnLifetime
	}
	if cfg.MaxConnIdleTime > 0 {
		conf.MaxConnIdleTime = cfg.MaxConnIdleTime
	}

	conf.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
