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
	User        User
}
type ProductWithOutUser struct {
	Id          uuid.UUID `gorm:"type:uuid" json:"Id"` // default:uuid_generate_v4() Must Create Feature UUID in DB: -> ostgres "CREATE EXTENSION IF NOT EXISTS "uuid-ossp";"
	Name        string    `json:"name" validate:"required"`
	Price       float64   `json:"price" validate:"required"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"create_at"`
	UpdatedAt   time.Time `json:"update_at"`
	UserId      uuid.UUID `json:"created_by"` // Hide in JSON format
}
type ProductWithUser struct {
	Id          uuid.UUID `gorm:"type:uuid" json:"Id"` // default:uuid_generate_v4() Must Create Feature UUID in DB: -> ostgres "CREATE EXTENSION IF NOT EXISTS "uuid-ossp";"
	Name        string    `json:"name" validate:"required"`
	Price       float64   `json:"price" validate:"required"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"create_at"`
	UpdatedAt   time.Time `json:"update_at"`
	UserId      uuid.UUID `json:"created_by"` // Hide in JSON format
	User        UserSimple
}

// Find All Products => Gorm.DB
func FindAllProducts() (*[]Product, error) {
	db := database.GetDB()
	products := &[]Product{}
	err := db.Model(&Product{}).Order("created_at desc").Find(products).Error
	return products, err
}

// Find Products Condition => Gorm.DB
func FindProductsByCondition(condition interface{}) (*[]ProductWithOutUser, error) {
	db := database.GetDB()
	products := &[]ProductWithOutUser{}
	err := db.Model(&Product{}).Order("created_at desc").Find(products).Error
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
	error := db.Model(&product).Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&User{}).Select("Id,email,address,created_at").Find(&UserSimple{})
	}).Where(byField).First(&product).Error
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
