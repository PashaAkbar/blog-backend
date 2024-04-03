package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	connection "github.com/pashaakbar/blog-backend/Connection"
	"github.com/pashaakbar/blog-backend/models"
)

func CreateAgenda(c *gin.Context) {
	var body struct {
		Title    string `json:"title"`
		Time     string `json:"time"`
		Location string `json:"location"`
	}

	if c.BindJSON(&body) != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}

	agenda := models.AgendaItem{
		Title:    body.Title,
		Time:     body.Time,
		Location: body.Location,
	}

	db, err := connection.ConnectToDB()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to connect to database",
		})
		return
	}
	if err := db.Create(&agenda).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to create agenda item",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Agenda item created",
		"agenda":  agenda,
	})
}
func GetAllAgenda(c *gin.Context) {
	var agenda []models.AgendaItem

	db, err := connection.ConnectToDB()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to connect to database",
		})
		return
	}
	if err := db.Find(&agenda).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to find agenda items",
		})
		return
	}
	c.JSON(200, gin.H{
		"agenda": agenda,
	})
}

func GetAgenda(c *gin.Context) {
	id := c.Param("id")

	var agenda models.AgendaItem

	db, err := connection.ConnectToDB()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to connect to database",
		})
		return
	}
	if err := db.Where("id = ?", id).First(&agenda).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to find agenda item",
		})
		return
	}
	c.JSON(200, gin.H{
		"agenda": agenda,
	})
}

func UpdateAgenda(c *gin.Context) {
	id := c.Param("id")
	var agenda models.AgendaItem
	// fmt.Println("id", id)

	var body struct {
		Title    string `json:"title"`
		Time     string `json:"time"`
		Location string `json:"location"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON format",
		})
		return
	}

	db, err := connection.ConnectToDB()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to connect to database",
		})
		return
	}

	if err := db.Model(&agenda).Where("id = ?", id).Updates(&body).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to update agenda item",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Agenda item updated",
		"agenda":  agenda,
	})
}

func DeleteAgenda(c *gin.Context) {
	id := c.Param("id")
	var agenda models.AgendaItem

	db, err := connection.ConnectToDB()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to connect to database",
		})
		return
	}

	if err := db.Where("id = ?", id).Delete(&agenda).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to delete agenda item",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Agenda item deleted",
	})
}
