package utils

import (
	"restfulAPI/Golang/config"
	"time"

	"github.com/golang-jwt/jwt"
)

var secret = "your-secret-key"

func GenerateAccessToken(username string, duration int64) (string, error) {
	config, err := config.InitConfig()
	if err != nil {
		return "", err
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Duration(duration)).Unix() // Token expires in 24 hours

	accessToken, err := token.SignedString([]byte(config.ServerConfig.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
