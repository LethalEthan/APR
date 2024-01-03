package auth

import (
	"APR/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// If the user is not logged in
func EnsureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Error if the token is empty or the user is not logged in
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if !loggedIn {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

// If the user is already logged in
func EnsureNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// No error if the token is not empty or the user is already logged in
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if loggedIn {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

// Set whether use is logged in or not
func SetUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie("token"); err == nil || token != "" {
			t, err := db.QueryTokenExists(token)
			if err != nil {
				log.Println("SUS: " + err.Error())
			}
			if t == token {
				c.Set("is_logged_in", true)
			} else {
				c.Set("is_logged_in", false)
			}
		} else {
			c.Set("is_logged_in", false)
		}
	}
}
