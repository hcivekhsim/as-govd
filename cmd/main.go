package main

import (
	"context"
	"log"
	"time"

	"github.com/hcivekhsim/as-govd/internal/database/postgres"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err := postgres.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("Ошибка init DB: %v", err) // Create custom error
	}
}

// Приложение запускается раньше, чем Postgres готов принимать подключения - FIX 