package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	connection "github.com/pashaakbar/blog-backend/Connection"
	"github.com/pashaakbar/blog-backend/controllers"
	"github.com/pashaakbar/blog-backend/initializer"
	"github.com/pashaakbar/blog-backend/middleware"
)

func init() {
	initializer.LoadEnv()

}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization") // Remove withCredentials from allowed headers
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")                                                             // Allow credentials
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

func main() {
	db, err := connection.ConnectToDB()

	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	// defer db.Close()

	connection.SyncDatabase(db)

	r := gin.Default()

	r.Use(corsMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/userData", middleware.RequireAuth, controllers.GetUserData3)
	r.POST("/createAgenda", middleware.RequireAuth, controllers.CreateAgenda)
	r.GET("/getAgenda/:id", middleware.RequireAuth, controllers.GetAgenda)
	r.GET("/getAllAgenda", middleware.RequireAuth, controllers.GetAllAgenda)
	r.PUT("/updateAgenda/:id", middleware.RequireAuth, controllers.UpdateAgenda)
	r.DELETE("/deleteAgenda/:id", middleware.RequireAuth, controllers.DeleteAgenda)

	r.Run()
}
