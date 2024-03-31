package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // Use v4 of golang-jwt library
	connection "github.com/pashaakbar/blog-backend/Connection"
	"github.com/pashaakbar/blog-backend/models"
)

func RequireAuth(c *gin.Context) {
	// Retrieve token from cookie
	tokenString, err := c.Cookie("Authentification")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil // Use your actual JWT secret from environment variables
	})
	if err != nil || !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Extract claims from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Check token expiry
	expiryTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().After(expiryTime) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Retrieve user from database based on token subject (user ID)
	db, err := connection.ConnectToDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to connect to database"})
		return
	}

	var user models.User
	if err := db.First(&user, "id = ?", claims["sub"]).Error; err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Set user in context
	c.Set("user", user)

	// Continue with the next middleware or route handler
	c.Next()
}
