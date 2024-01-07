package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
)

func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method != "POST" {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
			c.Abort()
			return
		}

		authHeader := c.Request.Header.Get("Authorization")
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		token := splitToken[1]
		if bcrypt.CompareHashAndPassword([]byte(os.Getenv("HASHED_TOKEN")), []byte(token)) != nil {
			// Log unauthorized access attempt here
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized",
				"message": "You are not authorized to access this resource",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
