package controllers

import (
	"golang_fp/database"
	"golang_fp/models"

	"github.com/gin-gonic/gin"
)

func DashboardHandler(c *gin.Context) {
	session, _ := models.Store.Get(c.Request, "session")
	username, ok := session.Values["username"].(string)

	if !ok || username == "" {
		c.Redirect(303, "/login")
		return
	}

	rows, err := database.DB.Query("SELECT id, hotelname, arrivaldate, departuredate, comment FROM bookings WHERE username = ?", username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error fetching bookings"})
		return
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		if err := rows.Scan(&booking.ID, &booking.Hotelname, &booking.Arrivaldate, &booking.Departuredate, &booking.Comment); err != nil {
			c.JSON(500, gin.H{"error": "Error scanning booking data"})
			return
		}
		bookings = append(bookings, booking)
	}

	successMessage := c.DefaultQuery("success", "")
	c.HTML(200, "dashboard.html", gin.H{
		"UserID":         username,
		"LogoutURL":      "/logout",
		"Bookings":       bookings,
		"NoBookings":     len(bookings) == 0,
		"SuccessMessage": successMessage,
	})
}