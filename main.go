package main

import (
	"restfulAPI/Golang/database"
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

	router.POST("/login", loginHandler)
	router.POST("/register", registerHandler)
	router.POST("/refresh", refreshHandler)
	router.POST("/forgot-password", forgotpasswordHander)
	router.POST("/reset-password/", resetpasswordHandler)
	router.GET("/protected", authMiddleware(), protectedHandler)
	router.Run(":8080")
}
