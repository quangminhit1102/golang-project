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
	User        User      `gorm:"foreignKey:UserId"`
}
type ProductWithOutUser struct {
	Id          uuid.UUID `gorm:"type:uuid" json:"Id"` // default:uuid_generate_v4() Must Create Feature UUID in DB: -> ostgres "CREATE EXTENSION IF NOT EXISTS "uuid-ossp";"
	Name        string    `json:"name" validate:"required"`
	Price       float64   `json:"price" validate:"required"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"create_at"`
	UpdatedAt   time.Time `json:"update_at"`
	UserId      uuid.UUID
}
type ProductWithUser struct {
	Id          uuid.UUID `gorm:"type:uuid" json:"Id"` // default:uuid_generate_v4() Must Create Feature UUID in DB: -> ostgres "CREATE EXTENSION IF NOT EXISTS "uuid-ossp";"
	Name        string    `json:"name" validate:"required"`
	Price       float64   `json:"price" validate:"required"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"create_at"`
	UpdatedAt   time.Time `json:"update_at"`
	UserId      uuid.UUID
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
	err := db.Model(&Product{}).Order("created_at desc").Where(condition).Find(products).Error
	return products, err
}

// Find Products Condition => Gorm.DB
func FindProductsWithPagination(
	page,
	pageSize int,
	searchString string,
	condition interface{},
) (*[]ProductWithOutUser, error) {

	db := database.GetDB()
	products := &[]ProductWithOutUser{}

	start := (page - 1) * pageSize

	err := db.Model(&Product{}).Limit(pageSize).Offset(start).Order("created_at desc").Where(condition).Where("name like ?", "%"+searchString+"%").Find(products).Error
	return products, err
}

// Create Product Function (Product -> uuid,error)
func CreateProduct(product *Product) (*Product, error) {
	db := database.GetDB()
	error := db.Create(product).Error
	return product, error
}

// Find Product Function (byfield -> product,error)
func FindOneProduct(byField interface{}) (*ProductWithUser, error) {
	db := database.GetDB()
	var product = &ProductWithUser{}
	error := db.Model(&Product{}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Model(&User{}).Find(&UserSimple{})
	}).Where(byField).First(product).Error
	return product, error
}

// Update Product Function (productId + updateField interface{} -> product, error)
func UpdateProduct(productId uuid.UUID, updateField interface{}) (*ProductWithUser, error) {
	db := database.GetDB()
	err := db.Find(&Product{}, productId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &ProductWithUser{}, err
	}
	error := db.Model(&Product{Id: productId}).Updates(updateField).Error

	product, err := FindOneProduct(productId)
	if err != nil {
		return &ProductWithUser{}, err
	}
	return product, error
}

// Update Product Function (productId + updateField interface{} -> product, error)
func DeleteProduct(productId uuid.UUID) (uuid.UUID, error) {
	db := database.GetDB()
	err := db.Delete(&Product{}, productId).Error
	return productId, err
}
