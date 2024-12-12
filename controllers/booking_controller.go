package controllers

import (
	"golang_fp/database"
	//"golang_fp/models"

	"github.com/gin-gonic/gin"
)

func BookHotelHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		// Retrieve the available hotels for booking
		rows, err := database.DB.Query("SELECT name FROM hotels")
		if err != nil {
			c.JSON(500, gin.H{"error": "Error fetching hotels"})
			return
		}
		defer rows.Close()

		var hotels []string
		for rows.Next() {
			var hotelname string
			if err := rows.Scan(&hotelname); err != nil {
				c.JSON(500, gin.H{"error": "Error scanning hotel data"})
				return
			}
			hotels = append(hotels, hotelname)
		}

		c.HTML(200, "book-hotel.html", gin.H{
			"Hotels": hotels,
		})
		return
	}

	if c.Request.Method == "POST" {
		username := c.PostForm("username")
		hotelname := c.PostForm("hotelname")
		arrivaldate := c.PostForm("arrivaldate")
		departuredate := c.PostForm("departuredate")
		comment := c.PostForm("comment")

		// Insert booking into the database
		_, err := database.DB.Exec("INSERT INTO bookings (username, hotelname, arrivaldate, departuredate, comment) VALUES (?, ?, ?, ?, ?)",
			username, hotelname, arrivaldate, departuredate, comment)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error booking hotel"})
			return
		}

		c.Redirect(303, "/dashboard?success=Booking%20successful")
	}
}

func DeleteBookingHandler(c *gin.Context) {
	id := c.Query("id") // Get the id from query parameters

	// Delete booking from the database
	_, err := database.DB.Exec("DELETE FROM bookings WHERE id = ?", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error deleting booking"})
		return
	}

	c.Redirect(303, "/dashboard?success=Booking%20deleted%20successfully")
}

func ModifyBookingHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		// Retrieve the booking ID from the query parameters
		bookingID := c.Query("id")

		// Fetch the booking details from the database
		row := database.DB.QueryRow("SELECT username, hotelname, arrivaldate, departuredate, comment FROM bookings WHERE id = ?", bookingID)
		var username, hotelname, arrivaldate, departuredate, comment string
		if err := row.Scan(&username, &hotelname, &arrivaldate, &departuredate, &comment); err != nil {
			c.JSON(500, gin.H{"error": "Error fetching booking details"})
			return
		}

		// Fetch the list of available hotels
		rows, err := database.DB.Query("SELECT name FROM hotels")
		if err != nil {
			c.JSON(500, gin.H{"error": "Error fetching hotels"})
			return
		}
		defer rows.Close()

		var hotels []string
		for rows.Next() {
			var hotel string
			if err := rows.Scan(&hotel); err != nil {
				c.JSON(500, gin.H{"error": "Error scanning hotel data"})
				return
			}
			hotels = append(hotels, hotel)
		}

		// Render the modification form
		c.HTML(200, "modify-booking.html", gin.H{
			"BookingID":            bookingID,
			"BookingUsername":      username,
			"BookingHotelname":     hotelname,
			"BookingArrivaldate":   arrivaldate,
			"BookingDeparturedate": departuredate,
			"Comment":              comment,
			"Hotels":               hotels,
		})
		return
	}

	if c.Request.Method == "POST" {
		// Retrieve form data
		bookingID := c.PostForm("id")
		username := c.PostForm("username")
		hotelname := c.PostForm("hotelname")
		arrivaldate := c.PostForm("arrivaldate")
		departuredate := c.PostForm("departuredate")
		comment := c.PostForm("comment")

		// Update the booking in the database
		_, err := database.DB.Exec(
			"UPDATE bookings SET username = ?, hotelname = ?, arrivaldate = ?, departuredate = ?, comment = ? WHERE id = ?",
			username, hotelname, arrivaldate, departuredate, comment, bookingID,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error updating booking"})
			return
		}

		c.Redirect(303, "/dashboard?success=Booking%20updated%20successfully")
	}
}

func UpdateBookingHandler(c *gin.Context) {
	id := c.PostForm("id")
	hotelname := c.PostForm("hotelname")
	arrivaldate := c.PostForm("arrivaldate")
	departuredate := c.PostForm("departuredate")
	comment := c.PostForm("comment")

	// Update booking in the database
	_, err := database.DB.Exec("UPDATE bookings SET hotelname = ?, arrivaldate = ?, departuredate = ?, comment = ? WHERE id = ?",
		hotelname, arrivaldate, departuredate, comment, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error updating booking"})
		return
	}

	c.Redirect(303, "/dashboard?success=Booking%20updated%20successfully")
}
