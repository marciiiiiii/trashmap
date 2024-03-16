package main

import (
	"trashmap_backend/internal/app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run()
}
