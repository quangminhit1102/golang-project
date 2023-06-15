package models

import (
	"errors"
	"restfulAPI/Golang/database"

	"gorm.io/gorm"
)

type User struct {
	Id                   int    `gorm:"type:uuid primary_key"`
	Email                string `gorm:"unique" json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8,password-strength"`
	Address              string
	Token                string
	RefreshToken         string
	ForgotPasswordToken  string
	ForgotPasswordExpire string
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

func SaveUser(user *User) (int, error) {
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
