package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) AddProductRouter(apiRouter *gin.RouterGroup) {
	productRouter := apiRouter.Group("product")

	productRouter.GET("/get-all", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "Get All"})
	})
	productRouter.POST("/add", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "Add Product"})
	})
	productRouter.PUT("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "Update Product"})
	})
	productRouter.DELETE("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": "Delete Product"})
	})
}
