package views

import (
	"api1/src/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/visitas", controllers.CreateVisitas)
	r.POST("/atracciones", controllers.CreateAtracciones)
	r.GET("/visitas/dia", controllers.GetVisitasFromDate)
	r.GET("/atracciones/dia", controllers.GetAtraccionesFromDate)
}