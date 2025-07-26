package main

import (
	"github.com/claudineijrdev/sib-crm-backend/internal/auth"
	"github.com/claudineijrdev/sib-crm-backend/internal/platform/database"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.Connect()

	r := gin.Default()

	api := r.Group("/api")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", auth.RegisterHandler)
		}
	}

	r.Run()
}