package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	mongodb "trashmap_backend/internal/app/mongoDB"
	"trashmap_backend/internal/app/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	loadEnv()

	databaseUri := os.Getenv("MONGODB_URI")
	if databaseUri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	dbHelper := mongodb.NewDatabaseHelper()

	dbHelper.Connect(databaseUri)

	r := gin.Default()
	routes.SetupRoutes(r, dbHelper)

	handler := cors.Default().Handler(r) // debug CORS policy

	// Cors policy for when specific origin can be set (flutter port is cahngiong with every restart)
	// handler := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://localhost:54653"},
	// 	AllowedMethods:   []string{"GET"},
	// 	AllowedHeaders:   []string{"Authorization"},
	// 	AllowCredentials: true,
	// }).Handler(r)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Try to gracefully shut down the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Disconnect from the database
	fmt.Print("Disconnecting from database...")
	databaseClient := dbHelper.GetClient()
	if err := databaseClient.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}

	if err := ctx.Err(); err != nil {
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

func loadEnv() {
	if err := godotenv.Load("internal/app/config/.env"); err != nil {
		log.Println("No .env file found")
	}
}
