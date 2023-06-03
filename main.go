package main

import (
	"fmt"
	"net/http"
	"restfulAPI/Golang/database"
	"restfulAPI/Golang/models"
	"restfulAPI/Golang/utils"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

var (
	router = gin.Default()
	secret = "your-secret-key"
)
// var validate *validator.Validate

func main() {
	db := database.Init()
	db.AutoMigrate(&models.User{})
	sqlDB, _ := db.DB()
	sqlDB.Close()

	router.POST("/login", loginHandler)
	router.POST("/refresh", refreshHandler)
	router.GET("/protected", authMiddleware(), protectedHandler)
	router.Run(":8080")
}

type LoginModel struct {
	//https://pkg.go.dev/github.com/go-playground/validator#hdr-Baked_In_Validators_and_Tags |GIN VALIDATOR TAG|
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,password-strength"`
}

func ValidatePassword(fl validator.FieldLevel) bool {
	return validatePassword(fl.Field().String())
}
func validatePassword(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func loginHandler(c *gin.Context) {
	var loginBody LoginModel
	c.Writer.Header().Set("Content-Type", "application/json")

	validate := validator.New()
	validate.RegisterValidation("password-strength", ValidatePassword)

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
	c.JSON(http.StatusOK, gin.H{"email": email, "password": password})

	// // You can perform your authentication logic here.
	// // For simplicity, let's assume the authentication is successful.

	// // Create the token
	// token := jwt.New(jwt.SigningMethodHS256)
	// claims := token.Claims.(jwt.MapClaims)
	// claims["username"] = email
	// claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	// // Sign the token with the secret key
	// tokenString, err := token.SignedString([]byte(secret))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func refreshHandler(c *gin.Context) {
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
	accessToken, err := generateAccessToken(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func generateAccessToken(username string) (string, error) {
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

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("username", claims["username"])
		c.Next()
	}
}

func protectedHandler(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{"message": "Protected endpoint", "username": username})
}
