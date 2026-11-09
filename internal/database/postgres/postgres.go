// Package postgres for connect db
package postgres

import (
	"context"
	"fmt"

	"github.com/hcivekhsim/as-govd/internal/config"
	"github.com/jackc/pgx/v5"
)

func ConnectDB(ctx context.Context) error {
	cfg := config.Load()
	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgresql://postgres:%s@db:%s/appeal_system?sslmode=disable", cfg.DBPassword, cfg.DBPort))
	if err != nil {
		return err // Create custom error
	}
	defer conn.Close(ctx)
	fmt.Println("Подключение к БД прошло успешно")
	return nil
}
