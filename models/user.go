package models

type User struct {
	Email                string
	Password             string
	Address              string
	Token                string
	ForgotPasswordToken  string
	ForgotPasswordExpire string
}
