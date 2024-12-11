package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

var (
	// Database instance
	db *sql.DB

	// Session store
	store = sessions.NewCookieStore([]byte("super-secret-key"))
)

type User struct {
	Username string
	Email    string
	Password string
}

type Hotel struct {
	Name        string
	Address     string
	Email       string
	Price       int
	Description string
	Rating      int
}

type Booking struct {
	ID            int
	Arrivaldate   time.Time
	Departuredate time.Time
	Hotelname     string
	Username      string
	Comment       string
}

func init() {
	var err error
	// Initialize the SQLite database with the name users.db
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	if db == nil {
		log.Fatalf("Failed to initialize database.")
	}

	// Create tables if they don't exist
	createTables()
}

func createTables() {
	// Create the "users" table if it doesn't exist
	createUserTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	`
	_, err := db.Exec(createUserTableQuery)
	if err != nil {
		log.Fatalf("Error creating users table: %v\n", err)
	}

	// Create the "hotels" table if it doesn't exist
	createHotelTableQuery := `
	CREATE TABLE IF NOT EXISTS hotels (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		address TEXT NOT NULL,
		email TEXT NOT NULL,
		price INTEGER NOT NULL,
		description TEXT NOT NULL,
		rating INTEGER NOT NULL
	);
	`
	_, err = db.Exec(createHotelTableQuery)
	if err != nil {
		log.Fatalf("Error creating hotels table: %v\n", err)
	}

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
	_, err = db.Exec(createBookingTableQuery)
	if err != nil {
		log.Fatalf("Error creating bookings table: %v\n", err)
	}
}

func main() {
	// Initialize Gin router
	router := gin.Default()
	// Load templates
	router.LoadHTMLGlob("templates/*")
	// Define routes
	router.GET("/", homeHandler)
	router.GET("/register", registerHandler)
	router.POST("/register", registerHandler)
	router.GET("/login", loginHandler)
	router.POST("/login", loginHandler)
	router.GET("/dashboard", dashboardHandler)
	router.GET("/logout", logoutHandler)
	router.GET("/add-hotel", addHotelHandler)
	router.POST("/add-hotel", addHotelHandler)
	router.GET("/book-hotel", bookHotelHandler)
	router.POST("/book-hotel", bookHotelHandler)
	//router.GET("/delete-booking/:id", deleteBookingHandler)
	router.GET("/delete-booking", deleteBookingHandler)
	router.GET("/search-hotel", searchHotelHandler)
	router.GET("/modify-booking", modifyBookingHandler)
	router.POST("/update-booking", updateBookingHandler)

	// Serve static files
	router.Static("/static", "./static")

	// Run the server
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(router.Run(":8080"))
}

// Render templates (adjusted for Gin)
func renderTemplate(c *gin.Context, tmpl string, data interface{}) {
	//err := c.HTML(200, tmpl, data)
	c.HTML(200, tmpl, data)
	//if err != nil {
	//	c.JSON(500, gin.H{"error": "Error rendering template"})
	//}
}

// Home handler
func homeHandler(c *gin.Context) {
	c.Redirect(303, "/login")
}

// Registration handler
func registerHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		errorMessage := c.DefaultQuery("error", "")
		successMessage := c.DefaultQuery("success", "")
		data := map[string]interface{}{
			"ErrorMessage":   errorMessage,
			"SuccessMessage": successMessage,
		}
		renderTemplate(c, "register.html", data)
		return
	}

	if c.Request.Method == "POST" {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")

		_, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
		if err != nil {
			c.Redirect(303, "/register?error=Username%20or%20email%20already%20exists")
			return
		}

		c.Redirect(303, "/login?success=Registration%20successful!%20You%20can%20now%20log%20in.")
	}
}

// Login handler
func loginHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		errorMessage := c.DefaultQuery("error", "")
		data := map[string]interface{}{
			"ErrorMessage": errorMessage,
		}
		renderTemplate(c, "login.html", data)
		return
	}

	if c.Request.Method == "POST" {
		username := c.PostForm("username")
		password := c.PostForm("password")

		var dbUsername, dbPassword string
		err := db.QueryRow("SELECT username, password FROM users WHERE username = ?", username).Scan(&dbUsername, &dbPassword)
		if err != nil || password != dbPassword {
			c.Redirect(303, "/login?error=Invalid%20username%20or%20password")
			return
		}

		// Create a session
		session, _ := store.Get(c.Request, "session")
		session.Values["username"] = username
		session.Save(c.Request, c.Writer)

		c.Redirect(303, "/dashboard")
	}
}

// Dashboard handler
func dashboardHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	username, ok := session.Values["username"].(string)

	if !ok || username == "" {
		c.Redirect(303, "/login")
		return
	}

	rows, err := db.Query("SELECT id, hotelname, arrivaldate, departuredate, comment FROM bookings WHERE username = ?", username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error fetching bookings"})
		return
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var booking Booking
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

// Logout handler
func logoutHandler(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	session.Values["username"] = nil
	session.Options.MaxAge = -1
	session.Save(c.Request, c.Writer)

	c.Redirect(303, "/login")
}

// Add hotel handler
func addHotelHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		renderTemplate(c, "add-hotel.html", nil)
		return
	}

	if c.Request.Method == "POST" {
		name := c.PostForm("name")
		address := c.PostForm("address")
		email := c.PostForm("email")
		price := c.PostForm("price")
		description := c.PostForm("description")
		rating := c.PostForm("rating")

		_, err := db.Exec("INSERT INTO hotels (name, address, email, price, description, rating) VALUES (?, ?, ?, ?, ?, ?)", name, address, email, price, description, rating)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error adding hotel"})
			return
		}

		c.Redirect(303, "/dashboard")
	}
}

// More handlers can be converted similarly...
// Book hotel handler
func bookHotelHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		// Retrieve the available hotels for booking
		rows, err := db.Query("SELECT name FROM hotels")
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

		renderTemplate(c, "book-hotel.html", gin.H{
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
		_, err := db.Exec("INSERT INTO bookings (username, hotelname, arrivaldate, departuredate, comment) VALUES (?, ?, ?, ?, ?)",
			username, hotelname, arrivaldate, departuredate, comment)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error booking hotel"})
			return
		}

		c.Redirect(303, "/dashboard?success=Booking%20successful")
	}
}

// Delete booking handler
func deleteBookingHandler(c *gin.Context) {
	id := c.Query("id") // Get the id from query parameters

	// Delete booking from the database
	_, err := db.Exec("DELETE FROM bookings WHERE id = ?", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error deleting booking"})
		return
	}

	c.Redirect(303, "/dashboard?success=Booking%20deleted%20successfully")
}

// Search hotel handler
func searchHotelHandler(c *gin.Context) {
	query := c.DefaultQuery("query", "")

	if query == "" {
		c.HTML(200, "search-hotel.html", gin.H{
			"ErrorMessage": "Please enter a hotel name to search.",
		})
		return
	}

	var hotel Hotel
	err := db.QueryRow("SELECT name, address, email, price, description, rating FROM hotels WHERE name LIKE ?", "%"+query+"%").
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

// Modify booking handler (GET)
func modifyBookingHandler(c *gin.Context) {
	id := c.DefaultQuery("id", "")

	// Retrieve booking details for modification
	var booking Booking
	err := db.QueryRow("SELECT id, username, hotelname, arrivaldate, departuredate, comment FROM bookings WHERE id = ?", id).
		Scan(&booking.ID, &booking.Username, &booking.Hotelname, &booking.Arrivaldate, &booking.Departuredate, &booking.Comment)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error fetching booking details"})
		return
	}

	renderTemplate(c, "modify-booking.html", gin.H{
		"Booking": booking,
	})
}

// Update booking handler (POST)
func updateBookingHandler(c *gin.Context) {
	id := c.PostForm("id")
	hotelname := c.PostForm("hotelname")
	arrivaldate := c.PostForm("arrivaldate")
	departuredate := c.PostForm("departuredate")
	comment := c.PostForm("comment")

	// Update booking in the database
	_, err := db.Exec("UPDATE bookings SET hotelname = ?, arrivaldate = ?, departuredate = ?, comment = ? WHERE id = ?",
		hotelname, arrivaldate, departuredate, comment, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error updating booking"})
		return
	}

	c.Redirect(303, "/dashboard?success=Booking%20updated%20successfully")
}
