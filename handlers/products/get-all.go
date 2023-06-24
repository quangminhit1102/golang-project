package handler

import (
	"fmt"
	"net/http"
	Product "restfulAPI/Golang/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var productsPerPage int = 2 // Products each page

// Get All Pruducts Handler
func GetAllProduct(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	searchString := c.DefaultQuery("search", "")
	userId := c.GetString("UserId")
	fmt.Printf(page, userId, searchString)
	products, err := Product.FindProductsByCondition(&Product.Product{UserId: uuid.Must(uuid.Parse(userId))})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	totalRecord := len(*products)
	c.JSON(http.StatusOK, gin.H{"totalRecord": totalRecord, "data": products})
}
