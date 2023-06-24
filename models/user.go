package models

import (
	"errors"
	"restfulAPI/Golang/database"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id                   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email                string    `gorm:"unique,Index" json:"email" validate:"required,email"`
	Password             string    `json:"password" validate:"required,min=8,password-strength"`
	Address              string
	Token                string
	RefreshToken         string
	ForgotPasswordToken  string
	ForgotPasswordExpire string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Products             []Product
}
type UserSimple struct {
	Id                   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email                string    `gorm:"unique,Index" json:"email" validate:"required,email"`
	Password             string    `json:"password" validate:"required,min=8,password-strength"`
	Address              string
	CreatedAt            time.Time
}

func FindOneByEmail(email string) (*User, error) {
	db := database.GetDB()
	var user = &User{}
	error := db.Where(&User{Email: email}).First(&user).Error
	return user, error
}

func FindOneByCondition(condition interface{}) (*User, error) {
	db := database.GetDB()
	var user = &User{}
	error := db.Where(condition).First(&user).Error
	return user, error
}

func SaveUser(user *User) (uuid.UUID, error) {
	db := database.GetDB()
	error := db.Create(user).Error
	return user.Id, error
}
func UpdateOneByEmail(email string, updateField interface{}) (*User, error) {
	db := database.GetDB()
	user, err := FindOneByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &User{}, err
	}
	error := db.Model(user).Updates(updateField).Error
	return user, error
}

// Golang Ignore Update The field that have nil, 0 Value
// Cmn mất cả buổi trời detect
func UpdateOneWithMap(email string, updateField map[string]interface{}) (*User, error) {
	db := database.GetDB()
	user, err := FindOneByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &User{}, err
	}
	error := db.Model(user).Updates(updateField).Error
	return user, error
}
