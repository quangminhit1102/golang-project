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

type ProductCreateRequest struct {
	Name        string  `json:"name" validate:"required"`
	Price       float64 `json:"price" validate:"gt=0"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
}
// Create Pruducts Handler
func CreateProduct(c *gin.Context) {
	productCreateReq := &ProductCreateRequest{}
	_ = c.ShouldBind(&productCreateReq)
	validate := validator.New()

	if err := validate.Struct(productCreateReq); err != nil {
		c.JSON(http.StatusOK, utils.NewValidatorError(err))
		return
	}
	username, _ := c.Get("username")
	user, err := User.FindOneByEmail(string(username))
	productId, err := Product.CreateProduct(
		&Product.Product{
			Id:          uuid.New(),
			Name:        productCreateReq.Name,
			Price:       productCreateReq.Price,
			Description: productCreateReq.Description,
			Image:       productCreateReq.Image,
			UserId:      uuid.Must(uuid.Parse("d67a4ed2-ed34-4fcc-954a-78349b267398")),
		})
	if err != nil {
		c.JSON(http.StatusOK, utils.NewValidatorError(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"Id": productId})
}
