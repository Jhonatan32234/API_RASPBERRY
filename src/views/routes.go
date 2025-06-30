package views

import (
	"api1/src/controllers"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	_ "api1/docs" // importa los docs generados
)

func RegisterRoutes(r *gin.Engine) {
	
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/visitas", controllers.CreateVisitas)
	r.POST("/atracciones", controllers.CreateAtracciones)
	r.GET("/visitas/dia", controllers.GetVisitasFromDate)
	r.GET("/atracciones/dia", controllers.GetAtraccionesFromDate)
}