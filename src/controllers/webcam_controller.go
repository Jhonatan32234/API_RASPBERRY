package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gocv.io/x/gocv"
)

var (
	cameraActive bool
	stopChan     chan bool
	mutex        sync.Mutex
)

func StartWebcam(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	if cameraActive {
		c.JSON(http.StatusOK, gin.H{"message": "La cámara ya está activa"})
		return
	}

	stopChan = make(chan bool)
	cameraActive = true

	go func() {
		webcam, err := gocv.OpenVideoCapture(0)
		if err != nil {
			log.Println("❌ Error al abrir la webcam:", err)
			cameraActive = false
			return
		}
		defer webcam.Close()

		img := gocv.NewMat()
		defer img.Close()

		log.Println("📷 Webcam iniciada")

		for {
			select {
			case <-stopChan:
				log.Println("🛑 Webcam detenida")
				return
			default:
				if ok := webcam.Read(&img); !ok || img.Empty() {
					continue
				}
				// Aquí va la lógica para contar personas (futura integración)
				fmt.Println("🧍 Frame capturado:", time.Now())
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Webcam iniciada"})
}

func StopWebcam(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	if !cameraActive {
		c.JSON(http.StatusOK, gin.H{"message": "La cámara no está activa"})
		return
	}

	stopChan <- true
	close(stopChan)
	cameraActive = false

	c.JSON(http.StatusOK, gin.H{"message": "Webcam detenida"})
}
