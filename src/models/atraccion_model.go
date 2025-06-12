package models

import (
	"api1/core/database"
	"api1/core/rabbitmq"
	"api1/src/entities"
	"encoding/json"
	"log"
	"os"
)

func GetAtraccionesFromDate(fecha string) ([]entities.Atraccion, error) {
	var atracciones []entities.Atraccion

	err := database.DB.Where("fecha >= ?", fecha).Find(&atracciones).Error
	if err != nil {
		return nil, err
	}

	if len(atracciones) == 0 {
		log.Println("No hay atracciones para enviar desde la fecha:", fecha)
		return atracciones, nil
	}

	if rabbitmq.Publish(atracciones, "atracciones_queue") {
		err = database.DB.Model(&entities.Atraccion{}).
			Where("fecha >= ? AND enviado = ?", fecha, false).
			Update("enviado", true).Error
		if err != nil {
			log.Println("Error actualizando atracciones como enviadas:", err)
		}
	}

	return atracciones, nil
}

func SaveAtracciones(input []entities.Atraccion) ([]entities.Atraccion, error) {
	var guardadas []entities.Atraccion

	for _, item := range input {
		item.Enviado = false
		if err := database.DB.Create(&item).Error; err != nil {
			log.Println("❌ Error al guardar atracción:", err)
			saveAtraccionToFile(item)
		} else {
			guardadas = append(guardadas, item)
		}
	}

	if len(guardadas) == 0 {
		log.Println("⚠️ Ninguna atracción fue guardada. No se enviará al broker.")
		return nil, nil
	}

	var toSend []entities.Atraccion
	database.DB.Where("enviado = ?", false).Find(&toSend)

	

	if len(toSend) > 0 && rabbitmq.Publish(toSend, "atracciones_queue") {
		database.DB.Model(&entities.Atraccion{}).Where("enviado = ?", false).Update("enviado", true)
	}

	return toSend, nil
}

func saveAtraccionToFile(data entities.Atraccion) {
	filePath := "core/database/saves/atracciones_saves.json"

	var prev []entities.Atraccion
	fileContent, _ := os.ReadFile(filePath)
	if len(fileContent) > 0 {
		_ = json.Unmarshal(fileContent, &prev)
	}

	prev = append(prev, data)

	content, err := json.MarshalIndent(prev, "", "  ")
	if err != nil {
		log.Println("❌ Error al serializar atracción para respaldo:", err)
		return
	}

	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		log.Println("❌ Error al escribir respaldo de atracción:", err)
	}
}
