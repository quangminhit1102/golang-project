package main

import (
	"restfulAPI/Golang/config"
	"restfulAPI/Golang/database"
	"restfulAPI/Golang/routers"

	Product "restfulAPI/Golang/models"
	User "restfulAPI/Golang/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var server = gin.Default()

// var validate *validator.Validate

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User.User{})
	db.AutoMigrate(&Product.Product{})
}

func main() {
	// Database init
	db := database.Init()
	Migrate(db)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Config Init
	config, err := config.InitConfig()
	if err != nil {
		return
	}

	// Router Init
	router := routers.Router{Server: server, Db: db}
	router.Init()


	server.Run(":" + config.ServerConfig.Port)
}
