package main

import (
	"github.com/claudineijrdev/sib-crm-backend/internal/container"
	"github.com/claudineijrdev/sib-crm-backend/internal/platform/database"
	"github.com/claudineijrdev/sib-crm-backend/internal/platform/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.Connect()

	// Criar container de dependências
	container := container.NewContainer(database.DB)

	r := gin.Default()

	// Adicionar middleware de telemetria global
	r.Use(middleware.TelemetryMiddleware(container.Telemetry))

	api := r.Group("/api")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", container.AuthHandler.Register)
			authRoutes.POST("/login", container.AuthHandler.Login)
		}
	}

	r.Run()
}