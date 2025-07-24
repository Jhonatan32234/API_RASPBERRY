package database

import (
	"api1/src/entities"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Connect() {
	/*
		127.0.0.1 <- ip jhonny
		192.168.1.245 <- ip alejandro
	*/
	dsn := "root:root@tcp(192.168.1.245:3307)/mydb?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal("Database connection failed", err)
	}

	db.AutoMigrate(&entities.Visitas{}, &entities.Atraccion{})
	DB = db
}
