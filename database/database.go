// database.go
package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Verify the connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	log.Println("Database connection established successfully!")
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

func CreateTables() {
	// Create the "users" table if it doesn't exist
	createUserTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	`
	_, err := DB.Exec(createUserTableQuery)
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
	_, err = DB.Exec(createHotelTableQuery)
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
	_, err = DB.Exec(createBookingTableQuery)
	if err != nil {
		log.Fatalf("Error creating bookings table: %v\n", err)
	}
}