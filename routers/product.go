package routers

import (
	productHandler "restfulAPI/Golang/handlers/products"

	"github.com/gin-gonic/gin"
)

func (r *Router) AddProductRouter(apiRouter *gin.RouterGroup) {
	productRouter := apiRouter.Group("product")

	productRouter.GET("/get-all", productHandler.GetAllProduct)
	productRouter.GET("/:id", productHandler.ProductDetail)
	productRouter.POST("/add", productHandler.CreateProduct)
	productRouter.PUT("/:id", productHandler.UpdateProduct)
	productRouter.DELETE("/:id", productHandler.DeleteProduct)
}
