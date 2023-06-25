package handler

import (
	"net/http"
	Product "restfulAPI/Golang/models"
	"restfulAPI/Golang/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type (
	ProductUpdateRequest struct {
		Name        string  `json:"name" validate:"required"`
		Price       float64 `json:"price" validate:"gt=0"`
		Description string  `json:"description"`
		Image       string  `json:"image"`
	}
	ProductUpdateResponse struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func UpdateProduct(c *gin.Context) {
	productUpdateReq := &ProductUpdateRequest{}
	productId := c.Param("id")
	UserId := c.GetString("UserId")
	productUUID, err := uuid.Parse(productId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Product not found!"})
		return
	}
	_ = c.ShouldBind(&productUpdateReq)
	validate := validator.New()

	if err := validate.Struct(productUpdateReq); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewValidatorError(err))
		return
	}

	product, err := Product.FindOneProduct(&Product.Product{Id: productUUID, UserId: uuid.Must(uuid.Parse(UserId))})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Product not found!"})
		return
	}

	productRes, err := Product.UpdateProduct(
		product.Id,
		&Product.Product{
			Name:        productUpdateReq.Name,
			Price:       productUpdateReq.Price,
			Description: productUpdateReq.Description,
			Image:       productUpdateReq.Image})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Internal Server Error!"})
		return
	}

	c.JSON(http.StatusNotFound, &ProductUpdateResponse{Success: true, Message: "Updated product!", Data: productRes})
}
