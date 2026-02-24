package psql

import (
	"context"
	"fmt"

	"github.com/OddKuru/core-accounts/internal/domain/vo"
	"github.com/OddKuru/core-accounts/internal/infra/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var (
	RoleValueFromDB = map[string]vo.AccountRoleType{
		"super_admin": vo.RoleSuperAdmin,
		"admin":       vo.RoleAdmin,
		"user":        vo.RoleUser,
	}

	RoleValueToDB = map[vo.AccountRoleType]string{
		vo.RoleUnknown:    "user",
		vo.RoleSuperAdmin: "super_admin",
		vo.RoleAdmin:      "admin",
		vo.RoleUser:       "user",
	}
)

type SQLQuerier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

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

	conf.HealthCheckPeriod = cfg.HealthCheckInterval

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
