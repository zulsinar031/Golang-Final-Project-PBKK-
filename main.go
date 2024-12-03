package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

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
	Arrivaldate   time.Time // Date field
	Departuredate time.Time // Date field
	Hotelname     string
	Username      string
	Comment       string
}

func init() {
	var err error
	// Initialize the SQLite database with the name users.db
	db, err = sql.Open("sqlite3", "C:/Users/leque/Golang-Final-Project-PBKK-/users.db")
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	if db == nil {
		log.Fatalf("Failed to initialize database.")
	}

	// Create the "users" table if it doesn't exist
	createUserTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	`
	_, err = db.Exec(createUserTableQuery)
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
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/add-hotel", addHotelHandler)
	http.HandleFunc("/book-hotel", bookHotelHandler)
	http.HandleFunc("/delete-booking/", deleteBookingHandler) // Use a dynamic URL for booking ID
	http.HandleFunc("/search-hotel", searchHotelHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Render templates
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(fmt.Sprintf("templates/%s", tmpl))
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

// Home handler
func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Registration handler
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Check if there's an error or success message in the URL query parameters
		errorMessage := r.URL.Query().Get("error")
		successMessage := r.URL.Query().Get("success")

		data := map[string]interface{}{
			"ErrorMessage":   errorMessage,
			"SuccessMessage": successMessage,
		}
		renderTemplate(w, "register.html", data)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		_, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
		if err != nil {
			// Redirect to registration page with an error message
			http.Redirect(w, r, "/register?error=Username%20or%20email%20already%20exists", http.StatusSeeOther)
			return
		}

		// Redirect to login page with a success message
		http.Redirect(w, r, "/login?success=Registration%20successful!%20You%20can%20now%20log%20in.", http.StatusSeeOther)
	}
}

// Login handler
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Check if there's an error message in the URL query parameters
		errorMessage := r.URL.Query().Get("error")
		data := map[string]interface{}{
			"ErrorMessage": errorMessage, // Pass the error message to the template
		}
		renderTemplate(w, "login.html", data)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var dbUsername, dbPassword string
		err := db.QueryRow("SELECT username, password FROM users WHERE username = ?", username).Scan(&dbUsername, &dbPassword)
		if err != nil || password != dbPassword {
			// Redirect to login page with error message
			http.Redirect(w, r, "/login?error=Invalid%20username%20or%20password", http.StatusSeeOther)
			return
		}

		// Create a session
		session, _ := store.Get(r, "session")
		session.Values["username"] = username
		session.Save(r, w)

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// Dashboard handler
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Get session and check if the user is logged in
	session, _ := store.Get(r, "session")
	username, ok := session.Values["username"].(string)

	// If the user is not logged in, redirect to the login page
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Query the database to get the bookings for the logged-in user
	rows, err := db.Query("SELECT id, hotelname, arrivaldate, departuredate, comment FROM bookings WHERE username = ?", username)
	if err != nil {
		http.Error(w, "Error fetching bookings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to store the bookings
	var bookings []Booking
	for rows.Next() {
		var booking Booking
		if err := rows.Scan(&booking.ID, &booking.Hotelname, &booking.Arrivaldate, &booking.Departuredate, &booking.Comment); err != nil {
			http.Error(w, "Error scanning booking data", http.StatusInternalServerError)
			return
		}
		bookings = append(bookings, booking)
	}

	successMessage := ""

	// Check if success parameter is in the URL
	if r.URL.Query().Get("success") == "true" {
		successMessage = "Booking deleted successfully!"
	}

	// Render the dashboard template, passing the success message and bookings data
	renderTemplate(w, "dashboard.html", map[string]interface{}{
		"UserID":         username,
		"LogoutURL":      "/logout",
		"Bookings":       bookings,
		"NoBookings":     len(bookings) == 0, // Add a flag indicating no bookings
		"SuccessMessage": successMessage,
	})
}

// Logout handler
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["username"] = nil
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Add hotel handler
func addHotelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Render the form to add a hotel
		renderTemplate(w, "add-hotel.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		// Get data from the form
		name := r.FormValue("name")
		address := r.FormValue("address")
		email := r.FormValue("email")
		price := r.FormValue("price")
		description := r.FormValue("description")
		rating := r.FormValue("rating") // Rating from hidden input

		// Insert into the hotels table
		_, err := db.Exec("INSERT INTO hotels (name, address, email, price, description, rating) VALUES (?, ?, ?, ?, ?, ?)", name, address, email, price, description, rating)
		if err != nil {
			http.Error(w, "Error adding hotel", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// Book hotel handler
func bookHotelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Get all hotels from the hotels table
		rows, err := db.Query("SELECT name FROM hotels")
		if err != nil {
			http.Error(w, "Error fetching hotels", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var hotels []string
		for rows.Next() {
			var hotelName string
			if err := rows.Scan(&hotelName); err != nil {
				http.Error(w, "Error scanning hotel name", http.StatusInternalServerError)
				return
			}
			hotels = append(hotels, hotelName)
		}

		// Render the form with the hotel options
		renderTemplate(w, "book-hotel.html", map[string]interface{}{
			"ErrorMessage": "",     // No error at this point
			"Hotels":       hotels, // Pass the hotels to the template
		})
		return
	}

	if r.Method == http.MethodPost {
		// Get data from the form
		arrivalDateStr := r.FormValue("arrivaldate")
		departureDateStr := r.FormValue("departuredate")
		hotelName := r.FormValue("hotelname")
		username := r.FormValue("username")
		comment := r.FormValue("comment")

		// Parse the dates
		arrivalDate, err := time.Parse("2006-01-02", arrivalDateStr)
		if err != nil {
			http.Error(w, "Invalid arrival date format", http.StatusBadRequest)
			return
		}
		departureDate, err := time.Parse("2006-01-02", departureDateStr)
		if err != nil {
			http.Error(w, "Invalid departure date format", http.StatusBadRequest)
			return
		}

		// Verify that the departure date is after the arrival date
		if departureDate.Before(arrivalDate) {
			renderTemplate(w, "book-hotel.html", map[string]interface{}{
				"ErrorMessage": "Departure date must be after the arrival date",
				"Hotels":       fetchHotels(),
			})
			return
		}

		// Verify that the username exists in the users table
		var userExists bool
		err = db.QueryRow("SELECT COUNT(1) FROM users WHERE username = ?", username).Scan(&userExists)
		if err != nil {
			http.Error(w, "Error checking username", http.StatusInternalServerError)
			return
		}
		if !userExists {
			renderTemplate(w, "book-hotel.html", map[string]interface{}{
				"ErrorMessage": "Username does not exist",
				"Hotels":       fetchHotels(),
			})
			return
		}

		// Verify that the hotel exists in the hotels table
		var hotelExists bool
		err = db.QueryRow("SELECT COUNT(1) FROM hotels WHERE name = ?", hotelName).Scan(&hotelExists)
		if err != nil {
			http.Error(w, "Error checking hotel", http.StatusInternalServerError)
			return
		}
		if !hotelExists {
			renderTemplate(w, "book-hotel.html", map[string]interface{}{
				"ErrorMessage": "Hotel does not exist",
				"Hotels":       fetchHotels(),
			})
			return
		}

		// Check if the hotel is available for the selected dates
		var overlapExists bool
		query := `
            SELECT COUNT(1) 
            FROM bookings 
            WHERE hotelname = ? 
            AND (
                (arrivaldate BETWEEN ? AND ?) OR 
                (departuredate BETWEEN ? AND ?)
            )
        `
		err = db.QueryRow(query, hotelName, arrivalDateStr, departureDateStr, arrivalDateStr, departureDateStr).Scan(&overlapExists)
		if err != nil {
			http.Error(w, "Error checking booking availability", http.StatusInternalServerError)
			return
		}

		if overlapExists {
			renderTemplate(w, "book-hotel.html", map[string]interface{}{
				"ErrorMessage": "This hotel is already booked for the selected dates",
				"Hotels":       fetchHotels(),
			})
			return
		}

		// Insert booking into the database
		_, err = db.Exec("INSERT INTO bookings (arrivaldate, departuredate, hotelname, username, comment) VALUES (?, ?, ?, ?, ?)", arrivalDateStr, departureDateStr, hotelName, username, comment)
		if err != nil {
			http.Error(w, "Error adding booking", http.StatusInternalServerError)
			return
		}

		// Redirect to the dashboard or a confirmation page
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

// Helper function to fetch hotels
func fetchHotels() []string {
	var hotels []string
	rows, err := db.Query("SELECT name FROM hotels")
	if err != nil {
		log.Println("Error fetching hotels:", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var hotelName string
		if err := rows.Scan(&hotelName); err != nil {
			log.Println("Error scanning hotel name:", err)
			continue
		}
		hotels = append(hotels, hotelName)
	}
	return hotels
}

// Delete booking handler
func deleteBookingHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the booking ID from the query parameter
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Booking ID is missing", http.StatusBadRequest)
		return
	}

	// Delete the booking from the database
	_, err := db.Exec("DELETE FROM bookings WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Error deleting booking", http.StatusInternalServerError)
		return
	}

	// Set a success message and redirect to the dashboard
	http.Redirect(w, r, "/dashboard?success=true", http.StatusSeeOther)
}

// Search hotel handler
func searchHotelHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		// If no query is provided, show an error message
		renderTemplate(w, "search-hotel.html", map[string]interface{}{
			"ErrorMessage": "Please enter a hotel name to search.",
		})
		return
	}

	// Query the database to search for hotels by name
	var hotel Hotel
	err := db.QueryRow("SELECT name, address, email, price, description, rating FROM hotels WHERE name LIKE ?", "%"+query+"%").
		Scan(&hotel.Name, &hotel.Address, &hotel.Email, &hotel.Price, &hotel.Description, &hotel.Rating)

	if err != nil {
		// If no hotel found or another error occurred, show an error message
		renderTemplate(w, "search-hotel.html", map[string]interface{}{
			"ErrorMessage": "No hotel found.",
		})
		return
	}

	// If a hotel is found, render the hotel information
	renderTemplate(w, "search-hotel.html", map[string]interface{}{
		"Hotel": hotel,
	})
}
