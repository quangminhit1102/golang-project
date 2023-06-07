package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/smtp"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// {"database": {"hello":"no such table", error: "not_exists"}}
type CommonError struct {
	Errors map[string]interface{} `json:"errors"`
}

// To handle the error returned by c.Bind in gin framework
// https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
func NewValidatorError(err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		// Switch Type Error And Custom Message
		switch v.Tag() {
		case "required":
			res.Errors[v.Field()] = fmt.Sprintf("The Field %s is required!", v.Field())
		case "email":
			res.Errors[v.Field()] = "Email is not valid!"
		case "required_without":
			res.Errors[v.Field()] = fmt.Sprintf("%s is required if %s is not supplied", v.Field(), v.Param())
		case "password-strength":
			res.Errors[v.Field()] = "Password mus contain at least 1 Uppercase, 1 Lower case, 1 number, 1 special character"
		case "min":
			res.Errors[v.Field()] = fmt.Sprintf("%s must be minimum %s characters length", v.Field(), v.Param())
		case "contains":
			res.Errors[v.Field()] = fmt.Sprintf("%s must contain at least one %s", v.Field(), v.Param())
		case "eqcsfield":
			res.Errors[v.Field()] = fmt.Sprintf("%s must match with %s", v.Field(), v.Param())
		case "lt", "ltfield":
			param := v.Param()
			if param == "" {
				param = time.Now().Format(time.RFC3339)
			}
			res.Errors[v.Field()] = fmt.Sprintf("%s must be less than %s", v.Field(), param)
		case "gt", "gtfield":
			param := v.Param()
			if param == "" {
				param = time.Now().Format(time.RFC3339)
			}
			res.Errors[v.Field()] = fmt.Sprintf("%s must be greater than %s", v.Field(), param)
		default:
			// if it's a tag for which we don't have a good format string yet we'll try using the default english translator
			english := en.New()
			translator := ut.New(english, english)
			if translatorInstance, found := translator.GetTranslator("en"); found {
				res.Errors[v.Field()] = v.Translate(translatorInstance)
			} else {
				res.Errors[v.Field()] = fmt.Errorf("%v", v).Error()
			}
		}

	}
	return res
}
func SendMail(toEmail, subject, body string) {
	// Choose auth method and set it up
	auth := smtp.PlainAuth("", "quangminhit.test01@gmail.com", "czkizmswbiljzmos", "smtp.gmail.com")
	// Here we do it all: connect to our server, set up a message and send it
	to := []string{toEmail}
	msg := []byte(body)
	err := smtp.SendMail("smtp.gmail.com:587", auth, "quangminhit.test01@gmail.com", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}

func RandomTokenGenerator() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
