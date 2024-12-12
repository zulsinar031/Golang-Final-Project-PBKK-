package controllers

import "github.com/gin-gonic/gin"

func HomeHandler(c *gin.Context) {
	c.Redirect(303, "/login")
}