package handlers

import (
	"errors"
	"fmt"
	"net/http"
	User "restfulAPI/Golang/models"
	"restfulAPI/Golang/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var secret = "your-secret-key"

type LoginModel struct {
	//https://pkg.go.dev/github.com/go-playground/validator#hdr-Baked_In_Validators_and_Tags |GIN VALIDATOR TAG|
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,password-strength"`
}

func RegisterHandler(c *gin.Context) {
	userRegister := &User.User{}
	_ = c.ShouldBind(&userRegister)
	validate := validator.New()
	validate.RegisterValidation("password-strength", utils.ValidatePassword)
	if err := validate.Struct(userRegister); err != nil {
		c.JSON(http.StatusOK, utils.NewValidatorError(err))
		return
	}
	email := userRegister.Email
	password := userRegister.Password

	_, err := User.FindOneByEmail(email)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// handle record not found
		c.JSON(http.StatusOK, gin.H{"error": true, "message": "Email Already exists!"})
		return
	}
	hashedByte, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	_, error := User.SaveUser(&User.User{Email: email, Password: string(hashedByte)})
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
	} else {
		c.JSON(http.StatusCreated, gin.H{"error": false, "message": "registed User Successfully!"})
	}

}
func LoginHandler(c *gin.Context) {
	var loginBody LoginModel
	c.Writer.Header().Set("Content-Type", "application/json")

	validate := validator.New()
	validate.RegisterValidation("password-strength", utils.ValidatePassword)

	// To Get From |Query| USING: c.DefaultQuery("<name>","<Default Value>")
	// To Get From |Param| USING: c.Param("name")
	// To Get From |JSON | USING: Bellow Code :)

	// if err := c.BindJSON(&loginBody); err != nil {
	// 	c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	// 	return
	// }
	_ = c.ShouldBind(&loginBody)
	if err := validate.Struct(loginBody); err != nil {
		c.JSON(http.StatusOK, utils.NewValidatorError(err))
		return
	}
	email := loginBody.Email
	password := loginBody.Password

	// c.JSON(http.StatusOK, gin.H{"email": email, "password": password})

	// Find User
	user, err := User.FindOneByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Email or Password!"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error!"})
		return
	}

	errorCompare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	// Check Password
	if user != nil && errorCompare == nil {
		// Create the token
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = email
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

		// Sign the token with the secret key
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"error": false, "message": "Authentication Successfully!", "token": tokenString})
	} else {
		c.JSON(http.StatusOK, gin.H{"error": true, "message": "Email or Password are invalid!"})
	}

}

func RefreshHandler(c *gin.Context) {
	tokenString := c.PostForm("refresh_token")

	// Validate the refresh token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Extract the username from the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	username := claims["username"].(string)

	// Generate a new access token
	accessToken, err := GenerateAccessToken(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func GenerateAccessToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func ProtectedHandler(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{"message": "Protected endpoint", "username": username})
}

func ForgotpasswordHander(c *gin.Context) {
	var jsonMap map[string]interface{}
	if err := c.ShouldBindJSON(&jsonMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	// Access the specific field
	email, ok := jsonMap["email"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field 'name'"})
		return
	}

	_, err := User.FindOneByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email doesn't exists!"})
		return
	}

	resetPasswordToken := utils.RandomTokenGenerator()
	body := "To: " + email + "\r\n" +
		"Subject: Reset Password for Application |Please follow these instruction|\r\n" +
		"\r\n" +
		"To reset password please click \r\n: " + "http://localhost:8080/reset-password/?email=" + email + "&token=" + resetPasswordToken
	utils.SendMail(email, "Reset Password for Application", body)
	User.UpdateOneByEmail(email, &User.User{ForgotPasswordToken: resetPasswordToken, ForgotPasswordExpire: "???"})
	c.JSON(http.StatusBadRequest, gin.H{"error": "Sent Mail To Reset Password Success!"})
}

type ResetPasswordReq struct {
	NewPassword     string `json:"newPassword" field:"New Password" validate:"required,min=8,password-strength"`
	ConfirmPassword string `json:"confirmPassword" field:"Confirm Password" validate:"required,min=8,password-strength,eqcsfield=NewPassword"`
}

func ResetpasswordHandler(c *gin.Context) {
	email := c.DefaultQuery("email", "")
	resetToken := c.DefaultQuery("token", "")
	if email == "" || resetToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request!"})
		return
	}
	user, err := User.FindOneByEmail(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request!"})
		return
	}
	// Check reset token
	if resetToken != user.ForgotPasswordToken {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token!"})
		return
	}
	// Bind Body
	resetPasswordReq := &ResetPasswordReq{}
	_ = c.ShouldBind(&resetPasswordReq)
	validate := validator.New()
	validate.RegisterValidation("password-strength", utils.ValidatePassword)

	if err := validate.Struct(resetPasswordReq); err != nil {
		c.JSON(http.StatusOK, utils.NewValidatorError(err))
		return
	}
	newPassword := resetPasswordReq.NewPassword
	hashedByte, _ := bcrypt.GenerateFromPassword([]byte(newPassword), 12)

	if _, err := User.UpdateOneByEmail(email, &User.User{Password: string(hashedByte)}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "messsage": "Reset password succesfully!"})
}
