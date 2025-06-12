package main

import (
    "api1/core/database"
	"api1/src/views"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	r := gin.Default()

	views.RegisterRoutes(r)
	port := "8080"
	r.Run(":" + port)
}
