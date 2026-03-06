package main

import (
	"go-mongo-project/db"
	"go-mongo-project/routes"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("No .env file found or error loading .env file")
	}

	db.Init()

	server := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	corsConfig := cors.DefaultConfig()

	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins != "" {
		rawOrigins := strings.Split(allowedOrigins, ",")
		origins := make([]string, 0, len(rawOrigins))
		for _, o := range rawOrigins {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				origins = append(origins, trimmed)
			}
		}
		corsConfig.AllowOrigins = origins
	} else {
		corsConfig.AllowOrigins = []string{"http://localhost:3007"}
	}

	corsConfig.AllowMethods = []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma", "Credential"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour

	server.Use(cors.New(corsConfig))

	routes.RegisterRoutes(server)

	server.Run(":8080")

}
