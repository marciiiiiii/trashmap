package main

import (
	"log"
	"os"
	"trashmap_backend/internal/app/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run()
}

func loadEnv() {
	if err := godotenv.Load("../../internal/app/config/.env"); err != nil {
		log.Println("No .env file found")
	}
}
