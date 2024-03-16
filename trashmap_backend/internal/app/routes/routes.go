package routes

import (
	"log"
	mongodb "trashmap_backend/internal/app/mongoDB"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, dbHelper *mongodb.DatabaseHelper) {

	r.GET("/ping", func(c *gin.Context) {
		result, err := dbHelper.FetchCollection("trashmap", "trashbins")
		if err != nil {
			log.Fatal("inside get: ", err)
		}
		c.JSON(200, result)
		// c.JSON(200, gin.H{
		// 	"message": "pong",
		// })
	})
}
