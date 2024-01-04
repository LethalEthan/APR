package routes

import (
	"APR/db"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	FirstNameEmpty = errors.New("The first name can't be empty")
	LastNameEmpty  = errors.New("The last name can't be empty")
	PasswordEmpty  = errors.New("The password can't be empty")
	EmailEmpty     = errors.New("The email isn't available")
)

var MemberList = make([]db.Member, 0, 10)

func init() {
	var err error
	MemberList, err = db.QueryAllMember()
	if err != nil {
		panic("Could not retrieve members in member table")
	}
}

// Check if the username and password combination is valid
func isUserValid(email, password string) bool {
	log.Println("UP: " + email + " " + password)
	for _, u := range MemberList {
		if u.Email == email && u.Password == password {
			return true
		}
	}
	return false
}

// Register a new user with the given username and password
func registerNewUser(fname, lname, email, password string) (*db.Member, error) {
	if strings.TrimSpace(fname) == "" {
		return nil, FirstNameEmpty
	}
	if strings.TrimSpace(lname) == "" {
		return nil, LastNameEmpty
	}
	if strings.TrimSpace(password) == "" {
		return nil, PasswordEmpty
	}
	if !isUserAvailable(email) {
		return nil, EmailEmpty
	}
	M := db.Member{FirstName: fname, LastName: lname, Email: email, Password: password}
	fmt.Println(fname + lname + email + password)
	MemberList = append(MemberList, M)
	db.InsertMember(M)
	return &M, nil
}

// Check if the supplied email is not being used
func isUserAvailable(email string) bool {
	for _, u := range MemberList {
		if u.Email == email {
			return false
		}
	}
	return true
}

func showProfile(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", nil)
}

func updatePersonalDetails(c *gin.Context) {
	// Obtain the form values by POST
	fname := c.PostForm("firstname")
	lname := c.PostForm("lastname")
	email := c.PostForm("email")
	password := c.PostForm("password")
	wellnessgoals := c.PostForm("wellnessgoals")
	dob := c.PostForm("date")
	fmt.Println("Form data:" + fname + lname + email + password + dob)
	if strings.TrimSpace(fname) == "" {
		c.HTML(http.StatusBadRequest, "personal-details.html", gin.H{
			"ErrorTitle":   "Update Failed",
			"ErrorMessage": "Firstname invalid"})
		return
	}
	if strings.TrimSpace(lname) == "" {
		c.HTML(http.StatusBadRequest, "personal-details.html", gin.H{
			"ErrorTitle":   "Update Failed",
			"ErrorMessage": "Lastname invalid"})
		return
	}
	if strings.TrimSpace(password) == "" {
		c.HTML(http.StatusBadRequest, "personal-details.html", gin.H{
			"ErrorTitle":   "Update Failed",
			"ErrorMessage": "Password invalid"})
		return
	}
	if strings.TrimSpace(email) == "" {
		c.HTML(http.StatusBadRequest, "personal-details.html", gin.H{
			"ErrorTitle":   "Update Failed",
			"ErrorMessage": "Email invalid"})
		return
	}
	if strings.TrimSpace(wellnessgoals) == "" {
		c.HTML(http.StatusBadRequest, "personal-details.html", gin.H{
			"ErrorTitle":   "Update Failed",
			"ErrorMessage": "WellnessGoals invalid"})
		return
	}
	if strings.TrimSpace(dob) == "" {
		c.HTML(http.StatusBadRequest, "personal-details.html", gin.H{
			"ErrorTitle":   "Update Failed",
			"ErrorMessage": "DOB invalid"})
		return
	}
	t, err := c.Cookie("token")
	if err != nil {
		c.HTML(http.StatusBadRequest, "personal-details.html", gin.H{
			"ErrorTitle":   "Update Failed",
			"ErrorMessage": "Token invalid"})
		return
	}
	if token, ok, err := db.QueryTokenExists(t); ok && err == nil && token == t {
		err = db.UpdateMember(db.Member{FirstName: fname, LastName: lname, Email: email, Password: password, WellnessGoals: wellnessgoals, DateOfBirth: dob})
		if err != nil {
			fmt.Println(err)
			c.HTML(http.StatusBadRequest, "personal-details.html", gin.H{
				"ErrorTitle":   "Update Failed",
				"ErrorMessage": err.Error()})
			return
		}
	} else {
		c.HTML(http.StatusBadRequest, "personal-details.html", gin.H{
			"ErrorTitle":   "Update Failed",
			"ErrorMessage": err.Error()})
		return
	}
}
