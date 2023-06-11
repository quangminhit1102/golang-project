package repositories

import User "restfulAPI/Golang/models"

type UserRepo interface {
	FindOneByEmail(email string) (*User.User, error)
	SaveUser(user *User.User) (int, error)
	UpdateOneByEmail(email string, updateField interface{}) (*User.User, error)
}
