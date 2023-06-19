package routers

import (
	"restfulAPI/Golang/handlers"
	"restfulAPI/Golang/middlewares"

	"github.com/gin-gonic/gin"
)

func (*Router) AddAuthenticationRouter(apiRouter *gin.RouterGroup) {
	authRouter := apiRouter.Group("auth")

	authRouter.POST("/login", handlers.LoginHandler)
	authRouter.POST("/register", handlers.RegisterHandler)
	authRouter.POST("/refresh", handlers.RefreshHandler)
	authRouter.POST("/forgot-password", handlers.ForgotpasswordHander)
	authRouter.POST("/reset-password/", handlers.ResetpasswordHandler)
	authRouter.GET("/protected", middlewares.AuthMiddleware(), handlers.ProtectedHandler)
}
