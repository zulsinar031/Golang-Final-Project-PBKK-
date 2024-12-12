# Golang CRUD Project - Hotel Booking Management

<p align="center">
    <img src="public/img/logo.png" alt="project_logo">
</p>

Welcome to the GitHub repository of the Golang CRUD application of Muhammad Izzul Sinar Mahadhika and Victor Lequeux Audran ! This project was created to manage Hotel Bookings, Users, and Hotels. It features functionality to create, modify, delete, and view bookings, along with a user-friendly interface for searching hotels.

Link to the video presentation:  
[https://youtu.be/7-1eEijIEaY?si=lxYNSDx8U1Jc9Krk](https://youtu.be/7-1eEijIEaY?si=lxYNSDx8U1Jc9Krk)

## Requirements

- `Go 1.21` or later
- `Gin Framework`
- `SQLite` (for local database)
- `HTML Templates` for UI

## Project Features

- **CRUD Operations** for Bookings
- **Hotel Search** functionality with search bar and results display
- **View content** of a single booking or hotel
- **Exception Handling** to prevent deletion of hotels linked to bookings
- **Gin Framework** for clean and efficient routing
- **HTML Templates** for responsive and user-friendly UI

## Database Structure

The project uses three primary tables:

- **Users**: Stores user details.
- **Hotels**: Stores hotel details including name and location.
- **Bookings**: Links users and hotels with details such as arrival date, departure date, and comments.

Here is a sample of the database schema for the `bookings` table:

```go
// Create the "bookings" table if it doesn't exist
	createBookingTableQuery := `
	CREATE TABLE IF NOT EXISTS bookings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		arrivaldate DATE NOT NULL,
		departuredate DATE NOT NULL,
		hotelname VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL,
		comment TEXT,
		FOREIGN KEY (hotelname) REFERENCES hotels(name),
		FOREIGN KEY (username) REFERENCES users(username)
	);
	`
```

<p align="center">
    <img src="static/img/capture_db_booking.png" alt="capture_db_booking">
</p>

## Controllers and Routes

- **auth_controller**: Handle registering, login in and login out 
- **booking_controller**: Manage booking a hotel, modifying and deleting a booking 
- **dashboard_controller**: Display  the dashboard with all the bookings of a person
- **home_controller**: Processes new bookings.

Here is an example of the `UpdateBookingHandler` function in `booking_controller.go`:

```go
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
```

## User Interface

The UI is built using HTML templates and includes:

- **Search Bar**: Allows users to search for hotels by name or location.
- **Booking Form**: Add or modify bookings with ease.
- **Responsive Design**: Buttons and links adapt for a seamless user experience.

<p align="center">
    <img src="static/img/capture_dashboard.png" alt="capture_dashboard">
</p>

Here is a preview of the hotel search results page:

<p align="center">
    <img src="static/img/capture_search_hotel.png" alt="capture_search_hotel">
</p>

## Exception Management

The project includes exception handling to ensure robust operations. For example:

- Prevents creating a user with an already taken email and username
- Notifies users of invalid search queries or missing booking details.

Here’s a code snippet to guarantee every username and email are unique :

```go
if c.Request.Method == "POST" {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")

		_, err := database.DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
		if err != nil {
			c.Redirect(303, "/register?error=Username%20or%20email%20already%20exists")
			return
		}

		c.Redirect(303, "/login?success=Registration%20successful!%20You%20can%20now%20log%20in.")
	}
```

<p align="center">
    <img src="static/img/capture_error.png" alt="error_handling">
</p>


