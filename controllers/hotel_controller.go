package controllers

import (
	"golang_fp/database"
	"golang_fp/models"

	"github.com/gin-gonic/gin"
)

func AddHotelHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(200, "add-hotel.html", nil)
		return
	}

	if c.Request.Method == "POST" {
		name := c.PostForm("name")
		address := c.PostForm("address")
		email := c.PostForm("email")
		price := c.PostForm("price")
		description := c.PostForm("description")
		rating := c.PostForm("rating")

		_, err := database.DB.Exec("INSERT INTO hotels (name, address, email, price, description, rating) VALUES (?, ?, ?, ?, ?, ?)", name, address, email, price, description, rating)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error adding hotel"})
			return
		}

		c.Redirect(303, "/dashboard")
	}
}

func SearchHotelHandler(c *gin.Context) {
	query := c.DefaultQuery("query", "")

	if query == "" {
		c.HTML(200, "search-hotel.html", gin.H{
			"ErrorMessage": "Please enter a hotel name to search.",
		})
		return
	}

	var hotel models.Hotel
	err := database.DB.QueryRow("SELECT name, address, email, price, description, rating FROM hotels WHERE name LIKE ?", "%"+query+"%").
		Scan(&hotel.Name, &hotel.Address, &hotel.Email, &hotel.Price, &hotel.Description, &hotel.Rating)

	if err != nil {
		// No hotel found or error occurred
		c.HTML(200, "search-hotel.html", gin.H{
			"ErrorMessage": "No hotel found.",
		})
		return
	}

	// If a hotel is found, pass the hotel data to the template
	c.HTML(200, "search-hotel.html", gin.H{
		"Hotel": hotel,
	})
}