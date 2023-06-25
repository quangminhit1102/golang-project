package handler

import (
	"fmt"
	"net/http"
	Product "restfulAPI/Golang/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var productsPerPage int = 2 // Products each page

type GetAllProductResponse struct {
	Success     bool        `json:"success"`
	Message     string      `json:"message"`
	TotalRecord int         `json:"totalRecord"`
	Data        interface{} `json:"data"`
}

// Get All Pruducts Handler
func GetAllProduct(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	searchString := c.DefaultQuery("search", "")
	userId := c.GetString("UserId")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 0 {
		pageInt = 1
	}
	fmt.Printf(page, userId, searchString, pageInt)

	products, err := Product.FindProductsWithPagination(pageInt, productsPerPage, searchString, &Product.Product{UserId: uuid.Must(uuid.Parse(userId))})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	totalRecord := len(*products)

	c.JSON(http.StatusOK, gin.H{"totalRecord": totalRecord, "data": products})
}
