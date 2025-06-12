package models

import (
	"api1/core/database"
	"api1/core/rabbitmq"
	"api1/src/entities"
	"encoding/json"
	"log"
	"os"
)

func GetVisitasFromDate(fecha string) ([]entities.Visitas, error) {
	var visitas []entities.Visitas

	// Obtener visitas desde la fecha que NO hayan sido enviadas (opcional)
	err := database.DB.Where("fecha >= ?", fecha).Find(&visitas).Error
	if err != nil {
		return nil, err
	}

	if len(visitas) == 0 {
		log.Println("No hay visitas para enviar desde la fecha:", fecha)
		return visitas, nil
	}

	// Publicar en RabbitMQ
	if rabbitmq.Publish(visitas, "visitas_queue") {
		// Marcar como enviadas
		err = database.DB.Model(&entities.Visitas{}).
			Where("fecha >= ? AND enviado = ?", fecha, false).
			Update("enviado", true).Error
		if err != nil {
			log.Println("Error actualizando visitas como enviadas:", err)
		}
	}

	return visitas, nil
}

func SaveVisitas(input []entities.Visitas) ([]entities.Visitas, error) {
	var guardadas []entities.Visitas

	for _, item := range input {
		item.Enviado = false
		if err := database.DB.Create(&item).Error; err != nil {
			log.Println("❌ Error al guardar visita:", err)
			saveVisitaToFile(item)
		} else {
			guardadas = append(guardadas, item)
		}
	}

	if len(guardadas) == 0 {
		log.Println("⚠️ Ninguna visita fue guardada. No se enviará al broker.")
		return nil, nil
	}

	var toSend []entities.Visitas
	database.DB.Where("enviado = ?", false).Find(&toSend)



	if len(toSend) > 0 && rabbitmq.Publish(toSend, "visitas_queue") {
		database.DB.Model(&entities.Visitas{}).Where("enviado = ?", false).Update("enviado", true)
	}

	return toSend, nil
}


func saveVisitaToFile(data entities.Visitas) {
	filePath := "core/database/saves/visitas_saves.json"

	var prev []entities.Visitas
	fileContent, _ := os.ReadFile(filePath)
	if len(fileContent) > 0 {
		_ = json.Unmarshal(fileContent, &prev)
	}

	prev = append(prev, data)

	content, err := json.MarshalIndent(prev, "", "  ")
	if err != nil {
		log.Println("❌ Error al serializar visita para respaldo:", err)
		return
	}

	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		log.Println("❌ Error al escribir respaldo de visita:", err)
	}
}
