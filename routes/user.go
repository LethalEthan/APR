package routes

import (
	"APR/db"
	"errors"
	"fmt"
	"log"
	"strings"
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
