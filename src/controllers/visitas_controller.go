package controllers

import (
	"api1/src/entities"
	"api1/src/models"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetVisitasFromDate(c *gin.Context) {
	fecha := c.Query("fecha") // Fecha en formato string (ej. "2025-06-10")
	if fecha == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere parámetro 'fecha'"})
		return
	}

	visitas, err := models.GetVisitasFromDate(fecha)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener visitas"})
		return
	}

	c.JSON(http.StatusOK, visitas)
}


func CreateVisitas(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo leer el cuerpo"})
		return
	}

	var unaVisita entities.Visitas
	if err := json.Unmarshal(body, &unaVisita); err == nil {
		result, err := models.SaveVisitas([]entities.Visitas{unaVisita})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar visita"})
			return
		}
		c.JSON(http.StatusCreated, result)
		return
	}

	var variasVisitas []entities.Visitas
	if err := json.Unmarshal(body, &variasVisitas); err == nil {
		result, err := models.SaveVisitas(variasVisitas)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar visitas"})
			return
		}
		c.JSON(http.StatusCreated, result)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
}
