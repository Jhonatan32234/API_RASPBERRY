package controllers

import (
	"log"
	"net/http"
	"os/exec"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	pythonCmd *exec.Cmd
	mu        sync.Mutex
)

// Iniciar detección (ejecutar script Python)
func StartWebcam(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	if pythonCmd != nil && pythonCmd.Process != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Detección ya en ejecución"})
		return
	}

	cmd := exec.Command("python3", "scripts/deteccion.py") // Ajusta ruta si es necesario

	err := cmd.Start()
	if err != nil {
		log.Printf("Error al iniciar script Python: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo iniciar la detección"})
		return
	}

	pythonCmd = cmd
	c.JSON(http.StatusOK, gin.H{"message": "Detección iniciada"})
}

// Detener detección (matar proceso Python)
func StopWebcam(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	if pythonCmd == nil || pythonCmd.Process == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No hay detección en ejecución"})
		return
	}

	err := pythonCmd.Process.Kill()
	if err != nil {
		log.Printf("Error al detener script Python: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo detener la detección"})
		return
	}

	pythonCmd = nil
	c.JSON(http.StatusOK, gin.H{"message": "Detección detenida"})
}
