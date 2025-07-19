// src/controllers/webcam_controller.go
package controllers

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

/* ---------- Tipo para la respuesta ---------- */

type peopleResp struct {
	People int `json:"people"`
}

/* ---------- Sesión global protegida ---------- */

// Cada puerta sólo puede tener un conteo activo a la vez.
// Si en el futuro hay varias puertas, usa un map[doorID]*session.
type session struct {
	cmd    *exec.Cmd
	stdout *bufio.Reader
}

var (
	mu   sync.Mutex
	curr *session
)

/* ---------- Auxiliares ---------- */

// inicia el script si no hay otro corriendo
func startCounting(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()
	if curr != nil {
		return errors.New("ya hay un conteo en curso")
	}

	// Ruta ABSOLUTA al script Python
	script := "/opt/museum/detect_people.py"

	cmd := exec.CommandContext(ctx, "python3", script, "--stream")
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return err
	}
	curr = &session{cmd: cmd, stdout: bufio.NewReader(stdout)}
	return nil
}

// detiene el script y devuelve la cuenta total
func stopCounting(ctx context.Context) (int, error) {
	mu.Lock()
	sess := curr
	curr = nil
	mu.Unlock()

	if sess == nil {
		return 0, errors.New("no hay conteo en curso")
	}

	// Señal suave: SIGINT permite que el script cierre limpiamente
	if err := sess.cmd.Process.Signal(os.Interrupt); err != nil {
		return 0, err
	}

	// Espera con timeout a que termine
	done := make(chan error, 1)
	go func() { done <- sess.cmd.Wait() }()
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case err := <-done:
		if err != nil {
			return 0, err
		}
	}

	// Lee la última línea JSON del stdout
	line, err := sess.stdout.ReadBytes('\n')
	if err != nil {
		return 0, err
	}
	var obj struct {
		People int `json:"people"`
	}
	if err := json.Unmarshal(line, &obj); err != nil {
		return 0, err
	}
	return obj.People, nil
}

/* ---------- Handlers ---------- */

// StartWebcam godoc
// @Summary Iniciar conteo de personas
// @Description Señal de puerta abierta: arranca el script Python para contar visitantes
// @Tags webcam
// @Produce json
// @Success 200 {object} map[string]string "status: counting"
// @Failure 409 {object} map[string]string "ya hay conteo"
// @Failure 500 {object} map[string]string "error interno"
// @Router /webcam/start [post]
func StartWebcam(c *gin.Context) {
	if err := startCounting(c.Request.Context()); err != nil {
		status := http.StatusConflict
		if !errors.Is(err, errors.New("ya hay un conteo en curso")) {
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "counting"})
}

// StopWebcam godoc
// @Summary Detener conteo y devolver total
// @Description Señal de puerta cerrada: detiene el script y responde con # de personas
// @Tags webcam
// @Produce json
// @Success 200 {object} peopleResp
// @Failure 409 {object} map[string]string "no hay conteo"
// @Failure 500 {object} map[string]string "error interno"
// @Router /webcam/stop [post]
func StopWebcam(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	total, err := stopCounting(ctx)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, errors.New("no hay conteo en curso")) {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Aquí podrías guardar en BD o publicar en RabbitMQ
	// models.SaveVisitas(...) o core.rabbitmq.Publisher.Publish(total)

	c.JSON(http.StatusOK, peopleResp{People: total})
}
