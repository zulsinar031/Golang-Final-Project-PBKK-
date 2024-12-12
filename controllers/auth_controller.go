package controllers

import (
	"golang_fp/database"
	"golang_fp/models"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		errorMessage := c.DefaultQuery("error", "")
		data := map[string]interface{} {
			"ErrorMessage": errorMessage,
		}
		c.HTML(200, "login.html", data)
		return
	}

	if c.Request.Method == "POST" {
		username := c.PostForm("username")
		password := c.PostForm("password")

		var dbUsername, dbPassword string
		err := database.DB.QueryRow("SELECT username, password FROM users WHERE username = ?", username).Scan(&dbUsername, &dbPassword)
		if err != nil || password != dbPassword {
			c.Redirect(303, "/login?error=Invalid%20username%20or%20password")
			return
		}

		// Create a session
		session, _ := models.Store.Get(c.Request, "session")
		session.Values["username"] = username
		session.Save(c.Request, c.Writer)

		c.Redirect(303, "/dashboard")
	}
}

func RegisterHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		errorMessage := c.DefaultQuery("error", "")
		successMessage := c.DefaultQuery("success", "")
		data := map[string]interface{}{
			"ErrorMessage":   errorMessage,
			"SuccessMessage": successMessage,
		}
		c.HTML(200, "register.html", data)
		return
	}

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
}

func LogoutHandler(c *gin.Context) {
	session, _ := models.Store.Get(c.Request, "session")
	session.Values["username"] = nil
	session.Options.MaxAge = -1
	session.Save(c.Request, c.Writer)

	c.Redirect(303, "/login")
}