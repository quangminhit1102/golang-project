package handler

import (
	"net/http"
	Product "restfulAPI/Golang/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Delete Pruduct Handler
func DeleteProduct(c *gin.Context) {
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
	if _, err := Product.DeleteProduct(product.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error when delete Product!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Delete Product successfully!"})
}
