package User

import (
	"restfulAPI/Golang/database"
	"restfulAPI/Golang/utils"
)

type User struct {
	Id                   int    `gorm:"primary_key,autoIncrement"`
	Email                string `gorm:"unique" json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8,password-strength"`
	Address              string
	Token                string
	ForgotPasswordToken  string
	ForgotPasswordExpire string
}

func FindOneByEmail(email string) (*User, error) {
	db := database.GetDB()
	var user = &User{}

	error := db.Where(&User{Email: email}).First(&user).Error
	return user, error
}
func SaveUser(user *User) (int, error) {
	db := database.GetDB()
	error := db.Create(user).Error
	return user.Id, error
}
func UpdateOneByEmail(email string) (*User, error) {
	db := database.GetDB()
	var user = &User{}
	error := db.Model(user).Updates(&User{ForgotPasswordToken: utils.RandomTokenGenerator(), ForgotPasswordExpire: "ExpireTime"}).Error
	return user, error
}
