package utils

import (
	"fmt"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// My own Error type that will help return my customized Error info
//
//	{"database": {"hello":"no such table", error: "not_exists"}}
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
		// can translate each error one at a time.
		//fmt.Println("gg",v.NameNamespace)
		// if v.Param() != "" {
		// 	res.Errors[v.Field()] = fmt.Sprintf("{%v: %v}", v.Tag(), v.Param())
		// } else {
		// 	res.Errors[v.Field()] = fmt.Sprintf("{key: %v}", v.Tag())
		// }
		switch v.Tag() {
		case "required":
			res.Errors[v.Field()] = fmt.Sprintf("The Field %s is required!", v.Field())
		case "email":
			res.Errors[v.Field()] = "Email is not valid!"
		case "required_without":
			res.Errors[v.Field()] = fmt.Sprintf("%s is required if %s is not supplied", v.Field(), v.Param())
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
