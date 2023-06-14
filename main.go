package main

import (
	"restfulAPI/Golang/config"
	"restfulAPI/Golang/database"
	"restfulAPI/Golang/handlers"
	"restfulAPI/Golang/middlewares"

	User "restfulAPI/Golang/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router = gin.Default()

// var validate *validator.Validate

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User.User{})
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

	router.POST("/login", handlers.LoginHandler)
	router.POST("/register", handlers.RegisterHandler)
	router.POST("/refresh", handlers.RefreshHandler)
	router.POST("/forgot-password", handlers.ForgotpasswordHander)
	router.POST("/reset-password/", handlers.ResetpasswordHandler)
	router.GET("/protected", middlewares.AuthMiddleware(), handlers.ProtectedHandler)
	router.Run(":" + config.ServerConfig.Port)
}
