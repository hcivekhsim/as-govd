package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcivekhsim/as-govd/internal/database/postgres"
	"github.com/hcivekhsim/as-govd/internal/tasks"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	pools, err := postgres.ConnectDB(ctx)
	if err != nil {
		log.Fatal(err) // Create custom error
	}
	defer pools.Close()

	userRepo := tasks.NewUserRepo(pools)
	authService := tasks.NewAuthService(userRepo)
	authHandle := tasks.NewAuthHandler(authService)

	r := gin.Default()
	api := r.Group("/api")
	{
		api.POST("register", authHandle.RegisterHandler)
		api.POST("login", authHandle.LoginHandler)
		api.GET("protected", authHandle.TestHandler)
	}

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
