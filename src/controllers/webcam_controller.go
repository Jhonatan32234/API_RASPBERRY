package controllers

import (
	"api1/src/entities"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	pythonCmd *exec.Cmd
	mu        sync.Mutex
)

// Iniciar detección (ejecutar script Python)
func StartWebcam(c *gin.Context) {
	log.Printf("entró al controller")

	logPath := "registro_personas.txt"
	err := os.WriteFile(logPath, []byte(""), 0644)
	if err != nil {
		log.Printf("Error al limpiar archivo de log: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if pythonCmd != nil && pythonCmd.Process != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Detección ya en ejecución"})
		return
	}

	cmd := exec.Command("./venv/bin/python", "./contador_visitas.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("despues de exec.comand")
	err = cmd.Start()
	if err != nil {
		log.Printf("Error al iniciar script Python: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo iniciar la detección"})
		return
	}

	pythonCmd = cmd
	c.JSON(http.StatusOK, gin.H{"message": "Detección iniciada"})
}

func StopWebcam(c *gin.Context) {
	log.Printf("entro al controller")
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

	logPath := "registro_personas.txt"
	data, err := os.ReadFile(logPath)
	log.Printf("Log leído:\n%s", string(data))

	if err != nil {
		log.Printf("Error al leer el log: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo leer el archivo de log"})
		return
	}

	lines := strings.Split(string(data), "\n")
	var totalIntervalo int
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		count, err := strconv.Atoi(line)
		if err == nil {
			totalIntervalo += count
		}
	}
	log.Printf("cantidad de ultima sesion: %v", totalIntervalo)

	visita := entities.Visitas{
		Visitantes: totalIntervalo,
		Hora:       time.Now().Format("15:04"),
		Fecha:      time.Now().Format("2006-01-02"),
		Zona:       "MZ",
		Enviado:    true,
	}

	log.Printf("ya se van a guardar las visitas")

	visitasJSON, err := json.Marshal([]entities.Visitas{visita})
	if err != nil {
		log.Printf("❌ Error al serializar visitas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo preparar la visita"})
		return
	}

	resp, err := http.Post("http://localhost:8080/visitas", "application/json", strings.NewReader(string(visitasJSON)))
	if err != nil {
		log.Printf("❌ Error al hacer POST a /visitas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar la visita"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("⚠️ Error en respuesta de /visitas: código %v", resp.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falló el guardado de la visita"})
		return
	}

	log.Printf("visitas guardadas... a mimir")
	c.JSON(http.StatusOK, gin.H{
		"message":    "Detección detenida",
		"visitantes": totalIntervalo,
	})
}
