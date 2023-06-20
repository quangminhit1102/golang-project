package routers

import (
	Auth "restfulAPI/Golang/handlers/auth"

	"github.com/gin-gonic/gin"
)

func (*Router) AddAuthenticationRouter(apiRouter *gin.RouterGroup) {
	authRouter := apiRouter.Group("auth")

	authRouter.POST("/login", Auth.LoginHandler)
	authRouter.POST("/register", Auth.RegisterHandler)
	authRouter.POST("/refresh", Auth.RefreshHandler)
	authRouter.POST("/forgot-password", Auth.ForgotpasswordHander)
	authRouter.POST("/reset-password/", Auth.ResetpasswordHandler)

	// authRouter.GET("/protected", middlewares.AuthMiddleware(), Auth.ProtectedHandler)
}
