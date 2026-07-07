package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ruangwali/internal/platform/config"
)

func OpenPostgres(
	ctx context.Context,
	cfg config.DatabaseConfig,
) (*pgxpool.Pool, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	poolConfig, err := pgxpool.ParseConfig(
		cfg.URL,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal parse konfigurasi PostgreSQL: %w",
			err,
		)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns

	poolConfig.MaxConnLifetime =
		cfg.MaxConnLifetime

	poolConfig.MaxConnIdleTime =
		cfg.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(
		ctx,
		poolConfig,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membuat PostgreSQL connection pool: %w",
			err,
		)
	}

	healthCtx, cancel := context.WithTimeout(
		ctx,
		cfg.HealthTimeout,
	)
	defer cancel()

	if err := pool.Ping(healthCtx); err != nil {
		pool.Close()

		return nil, fmt.Errorf(
			"gagal terhubung ke PostgreSQL: %w",
			err,
		)
	}

	return pool, nil
}
