package routes

import (
	"APR/db"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func showLoginPage(c *gin.Context) {
	Render(c, gin.H{
		"title": "Login",
	}, "login.html")
}

func showPersonalDetails(c *gin.Context) {
	token, err := c.Cookie("token") // Get token cookie
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided",
		})
		return
	}
	MID, err := db.QueryMemberIDByToken(token) // Get member ID by token
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided",
		})
		return
	}
	M, err := db.QueryMemberInfoByID(MID) // Get Member row by ID
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided",
		})
		return
	}
	Render(c, gin.H{
		"fname":         M.FirstName,
		"lname":         M.LastName,
		"email":         M.Email,
		"date":          M.DateOfBirth,
		"wellnessgoals": M.WellnessGoals,
	}, "personal-details.html")
}

func performLogin(c *gin.Context) {
	// Obtain the form values by POST
	email := c.PostForm("email")
	password := c.PostForm("password")
	log.Println("Form data:" + email + password)
	// Check if the username/password combination is valid
	if isUserValid(email, password) { // If the username/password is valid set the token in a cookie
		M, err := db.QueryMemberID(email, password)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"ErrorTitle":   "Login Failed",
				"ErrorMessage": "1" + err.Error(),
			})
			return
		}
		token := generateSessionToken()
		err = db.InsertToken(token, strconv.Itoa(M))
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"ErrorTitle":   "Login Failed",
				"ErrorMessage": "2" + err.Error(),
			})
			return
		}
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)
		Render(c, gin.H{
			"title": "Successful Login"}, "login-successful.html")

	} else {
		// If invalid show the error message on the login page
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided"})
	}
}

func generateSessionToken() string {
	return strconv.FormatInt(rand.Int63(), 16)
}

func logout(c *gin.Context) {
	t, err := c.Cookie("token")
	if err != nil {
		log.Println(err)
		c.SetCookie("token", "", -1, "", "", false, true)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	db.RevokeToken(t)
	c.SetCookie("token", "", -1, "", "", false, true) // Clear the cookie
	c.Redirect(http.StatusTemporaryRedirect, "/")     // Redirect to the home page
}

func showRegistrationPage(c *gin.Context) {
	Render(c, gin.H{"title": "Register"}, "register.html")
}

func ShowIndexPage(c *gin.Context) {
	Render(c, gin.H{"title": "Home Page"}, "index.html")
}

func register(c *gin.Context) {
	// Obtain the form values by POST
	fname := c.PostForm("firstname")
	lname := c.PostForm("lastname")
	email := c.PostForm("email")
	password := c.PostForm("password")
	log.Println("Form data:" + fname + lname + email + password)
	if _, err := registerNewUser(fname, lname, email, password); err == nil {
		// If the user is created, set the token in a cookie and log the user in
		M, err := db.QueryMemberID(email, password)
		if err != nil {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"ErrorTitle":   "Registration Failed",
				"ErrorMessage": err.Error(),
			})
			return
		}
		token := generateSessionToken()
		db.InsertToken(token, strconv.Itoa(M))
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)
		Render(c, gin.H{
			"title": "Successful registration & Login"},
			"login-successful.html")
	} else { // If error occured show on page
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"ErrorTitle":   "Registration Failed",
			"ErrorMessage": err.Error(),
		})
	}
}

func ShowNutritionPage(c *gin.Context) {
	c.HTML(http.StatusBadRequest, "nutrition.html", nil)
}

func ShowPinPage(c *gin.Context) {
	c.HTML(http.StatusBadRequest, "pin.html", nil)
}

func ShowCoachingPage(c *gin.Context) {
	c.HTML(http.StatusBadRequest, "coaching.html", nil)
}
