package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	connection "github.com/pashaakbar/blog-backend/Connection"
	"github.com/pashaakbar/blog-backend/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if c.BindJSON(&body) != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}

	// Hash Password

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed to hash password",
		})
		return
	}

	// Generate UUID for the user
	userUUID := uuid.New()

	// Create user with UUID
	user := models.User{
		ID:       userUUID.String(),
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hash),
	}
	db, err := connection.ConnectToDB()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to connect to database",
		})
		return
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to create user",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User created successfully",
		"user":    user,
	})

}
func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if c.BindJSON(&body) != nil {
		c.JSON(400, gin.H{
			"message": "failed to read body",
		})
		return
	}

	var user models.User

	db, err := connection.ConnectToDB()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to connect to database",
		})
		return
	}
	db.First(&user, "email = ?", body.Email)

	if user.ID == "" {
		c.JSON(404, gin.H{
			"message": "Invalid Email or Password",
		})
		return

	}

	// Compare Password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Invalid Email or Password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte("Sggdsjkbsas34bjkj432bbj"))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to generate token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authentification", tokenString, 3600*24, "", "", false, true)

	c.JSON(200, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}

func Validate(c *gin.Context) {

	user, _ := c.Get("user")
	c.JSON(200, gin.H{
		"message": "Im logged in",
		"user":    user,
	})
}

func Logout(c *gin.Context) {
	// Clear the authentication cookie
	c.SetCookie("Authentification", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

func GetUserData(c *gin.Context) {
	// Get the user ID from the URL parameter
	userID := c.Query("id")

	// Retrieve user data from the database based on the user ID
	var user models.User
	db, err := connection.ConnectToDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to connect to database",
		})
		return
	}
	db.First(&user, userID)

	if user.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	// Return user data
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func GetUserData2(c *gin.Context) {
	// Get the user ID from the token
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := claims["sub"].(string)

	// Retrieve user data from the database based on the user ID
	var user models.User
	db, err := connection.ConnectToDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to connect to database",
		})
		return
	}
	db.First(&user, userID)

	if user.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	// Return user data
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
