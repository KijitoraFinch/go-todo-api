package main

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if bcrypt.CompareHashAndPassword([]byte(os.Getenv("HASHED_TOKEN")), []byte(token)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized",
				"message": "This incident will be reported"})
			c.Abort()
		}
		c.Next()
	}
}
