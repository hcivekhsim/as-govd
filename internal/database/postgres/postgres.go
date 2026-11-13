// Package postgres for connect db
package postgres

import (
	"context"
	"fmt"

	"github.com/hcivekhsim/as-govd/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(ctx context.Context) (*pgxpool.Pool, error) {
	cfg := config.Load()
	pools, err := pgxpool.New(ctx, fmt.Sprintf("postgresql://postgres:%s@db:%s/appeal_system?sslmode=disable", cfg.DBPassword, cfg.DBPort))
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	if err := pools.Ping(ctx); err != nil {
		pools.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return pools, nil
}
