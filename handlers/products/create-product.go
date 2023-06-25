package handler

import (
	"net/http"
	Product "restfulAPI/Golang/models"
	User "restfulAPI/Golang/models"
	"restfulAPI/Golang/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type (
	ProductCreateRequest struct {
		Name        string  `json:"name" validate:"required"`
		Price       float64 `json:"price" validate:"gt=0"`
		Description string  `json:"description"`
		Image       string  `json:"image"`
	}
	ProductCreateResponse struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

// Create Pruducts Handler
func CreateProduct(c *gin.Context) {
	productCreateReq := &ProductCreateRequest{}
	_ = c.ShouldBind(&productCreateReq)
	validate := validator.New()

	if err := validate.Struct(productCreateReq); err != nil {
		c.JSON(http.StatusOK, utils.NewValidatorError(err))
		return
	}
	UserId := c.GetString("UserId")
	user, err := User.FindOneByCondition(&User.User{Id: uuid.Must(uuid.Parse(UserId))})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "User Not Found!"})
	}
	product, err := Product.CreateProduct(
		&Product.Product{
			Id:          uuid.New(),
			Name:        productCreateReq.Name,
			Price:       productCreateReq.Price,
			Description: productCreateReq.Description,
			Image:       productCreateReq.Image,
			UserId:      user.Id,
		})
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewValidatorError(err))
		return
	}
	c.JSON(http.StatusOK, &ProductCreateResponse{Success: true, Message: "Created Product Successfully!", Data: product})
}
