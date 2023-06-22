package models

import (
	"errors"
	"restfulAPI/Golang/database"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	Id          uuid.UUID `gorm:"type:uuid" json:"Id"` // default:uuid_generate_v4() Must Create Feature UUID in DB: -> ostgres "CREATE EXTENSION IF NOT EXISTS "uuid-ossp";"
	Name        string    `json:"name" validate:"required"`
	Price       float64   `json:"price" validate:"required"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"create_at"`
	UpdatedAt   time.Time `json:"update_at"`
	UserId      uuid.UUID `json:"-"` // Hide in JSON format
}

// Find All Products => Gorm.DB
func FindAllProduct() (*[]Product, error) {
	db := database.GetDB()
	products := &[]Product{}
	err := db.Order("created_at desc").Find(products).Error
	return products, err
}

// Create Product Function (Product -> uuid,error)
func CreateProduct(product *Product) (uuid.UUID, error) {
	db := database.GetDB()
	error := db.Create(product).Error
	return product.Id, error
}

// Find Product Function (byfield -> product,error)
func FindOneProduct(byField interface{}) (*Product, error) {
	db := database.GetDB()
	var product = &Product{}
	error := db.Where(byField).First(&product).Error
	return product, error
}

// Update Product Function (productId + updateField interface{} -> product, error)
func UpdateProduct(productId uuid.UUID, updateField interface{}) (*Product, error) {
	db := database.GetDB()
	product, err := FindOneProduct(&Product{Id: productId})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Product{}, err
	}
	error := db.Model(product).Updates(updateField).Error
	return product, error
}

// Update Product Function (productId + updateField interface{} -> product, error)
func DeleteProduct(productId uuid.UUID) (uuid.UUID, error) {
	db := database.GetDB()
	err := db.Delete(&Product{}, productId).Error
	return productId, err
}
