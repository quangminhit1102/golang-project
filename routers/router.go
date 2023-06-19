package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
	Server *gin.Engine
	Db     *gorm.DB
}

func (r *Router) Init() {
	apiRouter := r.Server.Group("v1")

	// Authentication router
	r.AddAuthenticationRouter(apiRouter)
	// Product Router
	r.AddProductRouter(apiRouter)
}
