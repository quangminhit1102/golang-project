package handler

import (
	"net/http"
	Product "restfulAPI/Golang/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Get Pruduct Detail Handler
func ProductDetail(c *gin.Context) {
	productId := c.Param("id")
	productUUID, err := uuid.Parse(productId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Product not found!"})
		return
	}
	product, err := Product.FindOneProduct(&Product.Product{Id: productUUID})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Product not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Get Data Successfully!", "data": product})
}
