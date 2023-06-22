package handler

import (
	"fmt"
	"net/http"
	Product "restfulAPI/Golang/models"

	"github.com/gin-gonic/gin"
)

var productsPerPage int = 2 // Products each page

// Get All Pruducts Handler
func GetAllProduct(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	username, _ := c.Get("username")
	fmt.Printf(page, username)
	products, err := Product.FindAllProduct()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	totalRecord := len(*products)
	c.JSON(http.StatusOK, gin.H{"totalRecord": totalRecord, "data": products})
}
