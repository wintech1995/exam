package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type TokenBody struct {
	Token string `json:"token" binding:"required"`
}

func JWTAuthen() gin.HandlerFunc {
	return func(c *gin.Context) {
		hmacSampleSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
		header := c.Request.Header.Get("Authorization")
		tokenString := strings.Replace(header, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return hmacSampleSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userId", claims["userId"])
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "forbidden", "message": err.Error()})
		}
		c.Next()
	}
}

func AuthenToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var json TokenBody
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hmacSampleSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
		tokenString := json.Token

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return hmacSampleSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userId", claims["userId"])
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "forbidden", "message": err.Error()})
		}
		c.Next()
	}
}
