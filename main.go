package main

import (
	"restfulAPI/Golang/database"
	"restfulAPI/Golang/handlers"
	"restfulAPI/Golang/middlewares"
	User "restfulAPI/Golang/models"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
	secret = "your-secret-key"
)

// var validate *validator.Validate

func main() {
	db := database.Init()
	db.AutoMigrate(&User.User{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	router.POST("/login", handlers.LoginHandler)
	router.POST("/register", handlers.RegisterHandler)
	router.POST("/refresh", handlers.RefreshHandler)
	router.POST("/forgot-password", handlers.ForgotpasswordHander)
	router.POST("/reset-password/", handlers.ResetpasswordHandler)
	router.GET("/protected", middlewares.AuthMiddleware(), handlers.ProtectedHandler)
	router.Run(":8080")
}
