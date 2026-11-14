package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcivekhsim/as-govd/internal/config"
	"github.com/hcivekhsim/as-govd/internal/database/postgres"
	"github.com/hcivekhsim/as-govd/internal/tasks"
)

func main() {
	cfg := config.Load()
	tasks.ApplyRuntimeConfig(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	pools, err := postgres.ConnectDB(ctx)
	if err != nil {
		log.Fatal(err) 
	}
	defer pools.Close()

	userRepo := tasks.NewUserRepo(pools)
	dictionaryRepo := tasks.NewDictionaryRepo(pools)
	appealRepo := tasks.NewAppealRepo(pools)
	authService := tasks.NewAuthService(userRepo)
	appealService := tasks.NewAppealService(appealRepo, dictionaryRepo)
	authHandle := tasks.NewAuthHandler(authService)
	dictionaryHandler := tasks.NewDictionaryHandler(dictionaryRepo)
	appealHandler := tasks.NewAppealHandler(appealService, appealRepo)

	r := gin.Default()
	r.Use(tasks.CORSMiddleware(tasksRuntimeOrigins(cfg)))
	api := r.Group("/api")
	{
		api.POST("register", authHandle.RegisterHandler)
		api.POST("login", authHandle.LoginHandler)
		api.POST("logout", authHandle.LogoutHandler)
		api.GET("protected", authHandle.TestHandler)

		dictionaries := api.Group("/dictionaries")
		{
			dictionaries.GET("/applicant-categories", dictionaryHandler.ApplicantCategories)
			dictionaries.GET("/topics", dictionaryHandler.Topics)
			dictionaries.GET("/appeal-kinds", dictionaryHandler.AppealKinds)
			dictionaries.GET("/appeal-types", dictionaryHandler.AppealTypes)
			dictionaries.GET("/physical-statuses", dictionaryHandler.PhysicalStatuses)
			dictionaries.GET("/social-categories", dictionaryHandler.SocialCategories)
			dictionaries.GET("/citizenships", dictionaryHandler.Citizenships)
			dictionaries.GET("/appeal-statuses", dictionaryHandler.AppealStatuses)
		}

		authed := api.Group("")
		authed.Use(tasks.AuthMiddleware(userRepo))
		{
			authed.PUT("/preferences", authHandle.PreferencesHandler)
			authed.POST("/appeals", appealHandler.Create)
			authed.GET("/appeals", appealHandler.ListMine)
			authed.GET("/appeals/:id", appealHandler.GetMine)
			authed.GET("/appeals/:id/files/:fileID/download", appealHandler.DownloadFileMine)
		}

		admin := api.Group("/admin")
		admin.Use(tasks.AuthMiddleware(userRepo), tasks.AdminMiddleware())
		{
			admin.GET("/appeals", appealHandler.ListAdmin)
			admin.GET("/appeals/:id", appealHandler.GetAdmin)
			admin.POST("/appeals/:id/status", appealHandler.ChangeStatus)
			admin.PATCH("/appeals/:id/status", appealHandler.ChangeStatus)
			admin.GET("/appeals/files/:fileID/download", appealHandler.DownloadFile)
		}
	}

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func tasksRuntimeOrigins(cfg *config.Config) []string {
	if cfg == nil || len(cfg.CORSAllowedOrigins) == 0 {
		return []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	}
	return cfg.CORSAllowedOrigins
}
