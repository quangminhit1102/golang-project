package middlewares

import (
	"fmt"
	"net/http"
	"restfulAPI/Golang/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Config Init
		config, err := config.InitConfig()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Init Config Error!"})
			c.Abort()
			return
		}

		// Get Token from Header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		// Extract Token
		if strings.Contains(tokenString, "Bear") {
			tokenString = strings.Split(tokenString, " ")[1] // Get second Item of string array
		}

		// Parse Token - Check Valid token With JwtSecretKey
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(config.ServerConfig.JwtSecretKey), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Parse Claim from Token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("username", claims["username"])
		c.Set("UserId",)

		c.Next() // Next Handler
	}
}
