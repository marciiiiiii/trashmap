package main

import (
	"log"
	mongodb "trashmap_backend/internal/app/mongoDB"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, dbHelper *mongodb.DatabaseHelper) {

	r.GET("/trashbins", verifyJWT(func(c *gin.Context) {
		getTrashbins(c, dbHelper)
	}))
}

func getTrashbins(c *gin.Context, dbHelper *mongodb.DatabaseHelper) {
	result, err := dbHelper.FetchCollection("trashmap", "trashbins")
	if err != nil {
		log.Fatal("inside get: ", err)
	}
	c.JSON(200, result)
}
