package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() *gorm.DB {
	dsn := "host=localhost port=5432 dbname=golang user=postgres password=123456 sslmode=prefer connect_timeout=10 TimeZone=Asia/Ho_Chi_Minh"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("Failed to connect DB:", err)
	}

	DB = db
	return DB
}
func GetDB() *gorm.DB {
	return DB
}
