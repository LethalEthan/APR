package auth

import (
	"APR/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ensure user is logged in and abort if not
func EnsureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Error if the token is empty or the user is not logged in
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool) // type convert any to bool
		if !loggedIn {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"ErrorTitle":   "Login Failed",
				"ErrorMessage": "Invalid credentials provided"})
			c.Abort()
			// c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

// Ensure user is logged in and authenticate with token
func EnsureAuthenticatedLoggedIn(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is logged in
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if !loggedIn { // if not abort
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"ErrorTitle":   "Login Failed",
				"ErrorMessage": "Invalid credentials provided"})
			c.Abort()
			return
		}
		t, err := c.Cookie("token") // Get cookie
		if err != nil {             // abort on error and set logged in to false
			log.Println(err)
			c.Set("is_logged_in", false)
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"ErrorTitle":   "Login Failed",
				"ErrorMessage": "Invalid credentials provided"})
			c.Abort()
			return
		}
		if token, ok, err := db.QueryTokenExists(t); ok && err != nil { // if user is logged in check token exists in DB
			if t != token { // if token DB != cookie token abort and set logged in false
				log.Println("Token mismatch!")
				c.Set("is_logged_in", false)
				c.HTML(http.StatusBadRequest, "login.html", gin.H{
					"ErrorTitle":   "Login Failed",
					"ErrorMessage": "Invalid credentials provided"})
				c.Abort()
				return
			}
			return
		} else {
			c.Set("is_logged_in", false)
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"ErrorTitle":   "Login Failed",
				"ErrorMessage": "Invalid credentials provided"})
			c.Abort()
			return
		}
	}
}

// Ensure user is not logged in
func EnsureNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// No error if the token is not empty or the user is already logged in
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if loggedIn { // abort if logged in
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

// Set whether use is logged in or not
func SetUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie("token"); err == nil || token != "" {
			t, _, err := db.QueryTokenExists(token)
			if err != nil {
				log.Println("SUS: " + err.Error())
			}
			if t == token { // Check DB token against cookie token yum yum
				c.Set("is_logged_in", true)
			} else {
				c.Set("is_logged_in", false)
			}
		} else {
			c.Set("is_logged_in", false)
		}
	}
}
