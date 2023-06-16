package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"restfulAPI/Golang/config"
	User "restfulAPI/Golang/models"
	"restfulAPI/Golang/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginModel struct {
	//https://pkg.go.dev/github.com/go-playground/validator#hdr-Baked_In_Validators_and_Tags |GIN VALIDATOR TAG|
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ResetPasswordReq struct {
	NewPassword     string `json:"newPassword" field:"New Password" validate:"required,min=8,password-strength"`
	ConfirmPassword string `json:"confirmPassword" field:"Confirm Password" validate:"required,min=8,password-strength,eqcsfield=NewPassword"`
}

func RegisterHandler(c *gin.Context) {
	// New Model
	userRegister := &User.User{}
	// Bind Model
	_ = c.ShouldBind(&userRegister)
	// Validate
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
		// Handle record not found
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Email Already exists!"})
		return
	}
	hashedByte, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	_, error := User.SaveUser(&User.User{Email: email, Password: string(hashedByte)})
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Internal Server Error!"})
	} else {
		c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Registered User Successfully!"})
	}
}
func LoginHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	// Init config
	config, err := config.InitConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error!"})
	}
	var loginBody LoginModel // New Login Model

	// Validator
	validate := validator.New()

	// To Get From |Query| USING: c.DefaultQuery("<name>","<Default Value>") ===================
	// To Get From |Param| USING: c.Param("name") ==============================================
	// To Get From |JSON | USING: Bellow Code :) ===============================================

	// if err := c.BindJSON(&loginBody); err != nil {
	// 	c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	// 	return
	// }
	_ = c.ShouldBind(&loginBody) // Bind Model
	// Validate
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
		tokenString, err := utils.GenerateAccessToken(email, int64(time.Minute))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		// Create Refresh Token
		refresh_token, err := utils.GenerateAccessToken(email, int64(time.Hour*24))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
			return
		}

		_, saveUserErr := User.UpdateOneByEmail(email, &User.User{Token: tokenString, RefreshToken: refresh_token})
		if saveUserErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error!"})
			return
		}

		c.SetCookie("token", tokenString, config.ServerConfig.AccessTokenMaxAge, "/", "localhost", false, true)
		c.SetCookie("refresh_token", refresh_token, config.ServerConfig.RefreshTokenMaxAge, "/", "localhost", false, true)
		c.SetCookie("logged_in", "true", config.ServerConfig.AccessTokenMaxAge, "/", "localhost", false, false)

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Authentication Successfully!", "token": tokenString, "refresh_token": refresh_token})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Email or Password are invalid!"})
	}
}

func RefreshHandler(c *gin.Context) {
	// Init Config
	config, err := config.InitConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error!"})
	}
	// Binding Refresh Token
	var jsonMap map[string]interface{}
	if err := c.ShouldBindJSON(&jsonMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request!"})
		return
	}
	// Find User By Refresh Token

	// Access the specific field
	refreshToken, ok := jsonMap["refreshToken"].(string)
	if !ok && strings.Trim(refreshToken, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required!"})
		return
	}

	// Find User By refresh Token
	user, err := User.FindOneByCondition(&User.User{RefreshToken: refreshToken})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh Token!"})
		return
	}

	// Parse Token
	if strings.Contains(refreshToken, "Bear") {
		refreshToken = strings.Split(refreshToken, " ")[1] // Get second Item of string array
	}
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(config.ServerConfig.JwtSecretKey), nil
	})
	// Validate the refresh token
	if err != nil || !token.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh Token!"})
		return
	}

	// Extract the username from the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	email := claims["username"].(string)

	// Generate a new access token
	newToken, err := utils.GenerateAccessToken(email, int64(time.Minute))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}
	// Generate a new refresh token
	newRefreshToken, err := utils.GenerateAccessToken(email, int64(time.Hour*24))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}
	// Update User
	User.UpdateOneByEmail(user.Email, User.User{RefreshToken: newRefreshToken, Token: newToken})

	// Set Cookie
	c.SetCookie("token", newToken, config.ServerConfig.AccessTokenMaxAge, "/", "localhost", false, true)
	c.SetCookie("refresh_token", newRefreshToken, config.ServerConfig.RefreshTokenMaxAge, "/", "localhost", false, true)
	c.SetCookie("logged_in", "true", config.ServerConfig.AccessTokenMaxAge, "/", "localhost", false, false)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Re-create Token successfully!", "token": newToken, "refresh_token": newRefreshToken})
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
	User.UpdateOneByEmail(email, &User.User{ForgotPasswordToken: resetPasswordToken, ForgotPasswordExpire: time.Now().Format("2006-01-02 15:04:05")})
	c.JSON(http.StatusOK, gin.H{"message": "Sent Mail To Reset Password Success!"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	User.UpdateOneWithMap(email, map[string]interface{}{"ForgotPasswordToken": "", "ForgotPasswordExpire": ""})
	// c.JSON(http.StatusOK, gin.H{"success": true, "message": "Reset password successfully!"})
	ResponseHandler(c, true, "Reset password successfully!", http.StatusOK)
}

func ResponseHandler(c *gin.Context, success bool, message string, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	c.JSON(statusCode, gin.H{"error": success, "message": message})
}
