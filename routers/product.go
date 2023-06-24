package routers

import (
	productHandler "restfulAPI/Golang/handlers/products"
	middlewares "restfulAPI/Golang/middlewares"

	"github.com/gin-gonic/gin"
)

func (r *Router) AddProductRouter(apiRouter *gin.RouterGroup) {
	productRouter := apiRouter.Group("product")

	productRouter.GET("/get-all", middlewares.AuthMiddleware(), productHandler.GetAllProduct)
	productRouter.GET("/:id", middlewares.AuthMiddleware(), productHandler.ProductDetail)
	productRouter.POST("/add", middlewares.AuthMiddleware(), productHandler.CreateProduct)
	productRouter.PUT("/:id", middlewares.AuthMiddleware(), productHandler.UpdateProduct)
	productRouter.DELETE("/:id", middlewares.AuthMiddleware(), productHandler.DeleteProduct)
}
