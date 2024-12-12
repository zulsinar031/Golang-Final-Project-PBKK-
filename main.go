package main

import (
	"fmt"
	"golang_fp/controllers"
	"golang_fp/database"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize database
	database.Init()
	// Initialize Gin router
	router := gin.Default()
	// Load templates
	router.LoadHTMLGlob("templates/*")
	// Serve static files
	router.Static("/static", "./static")

	// Define routes
	router.GET("/", controllers.HomeHandler)
	router.GET("/register", controllers.RegisterHandler)
	router.POST("/register", controllers.RegisterHandler)
	router.GET("/login", controllers.LoginHandler)
	router.POST("/login", controllers.LoginHandler)
	router.GET("/dashboard", controllers.DashboardHandler)
	router.GET("/logout", controllers.LogoutHandler)
	router.GET("/add-hotel", controllers.AddHotelHandler)
	router.POST("/add-hotel", controllers.AddHotelHandler)
	router.GET("/book-hotel", controllers.BookHotelHandler)
	router.POST("/book-hotel", controllers.BookHotelHandler)
	router.GET("/delete-booking", controllers.DeleteBookingHandler)
	router.GET("/search-hotel", controllers.SearchHotelHandler)
	router.GET("/modify-booking", controllers.ModifyBookingHandler)
	router.POST("/update-booking", controllers.UpdateBookingHandler)

	// Run the server
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(router.Run(":8080"))
}