package database

import (
	"log"
	"restfulAPI/Golang/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init Database
func Init() *gorm.DB {
	// Config Init
	config, err := config.InitConfig()
	if err != nil {
		log.Fatalln("Failed to Get Config:", err)
	}
	//// PostGresDB==
	dsn := config.ConnectionString

	//// MySQL ======
	//// dsn := "root:123456@tcp(localhost:3306)/golang?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect DB:", err)
	}
	DB = db

	return DB
}

// Get Database Instance
func GetDB() *gorm.DB {
	return DB
}
