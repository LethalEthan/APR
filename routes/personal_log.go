package routes

import (
	"APR/db"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var personalLogList []db.PersonalLog

func init() {
	var err error
	personalLogList, err = db.QueryPersonalLogs()
	if err != nil {
		panic(err)
		return
	}
}

// Return a list of all the articles
func GetAllLogsByMemberID(MemberID int) []db.PersonalLog {
	memberLogs := make([]db.PersonalLog, 0, 10)
	for _, p := range personalLogList {
		if p.MemberID == MemberID {
			memberLogs = append(memberLogs, p)
		}
	}
	return memberLogs
}

// Fetch an article based on the ID supplied
func GetLogByID(LogID int, MemberID int) (*db.PersonalLog, error) {
	for _, p := range personalLogList {
		if p.LogID == LogID && p.MemberID == MemberID {
			return &p, nil
		}
	}
	return nil, errors.New("Log not found")
}

func ShowLogs(c *gin.Context) {
	t, err := c.Cookie("token")
	if err != nil { // abort request on error and show login page
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided"})
	}
	if token, ok, err := db.QueryTokenExists(t); ok && err == nil { // if user is logged in check token exists in DB
		if t != token { // if token DB != cookie token abort and set logged in false
			log.Println("Token mismatch!")
			c.Set("is_logged_in", false)
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"ErrorTitle":   "Login Failed",
				"ErrorMessage": "Invalid credentials provided"})
			return
		}
	} else {
		log.Println(err)
		log.Println("Token not found in DB!")
		c.Set("is_logged_in", false)
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided"})
		return
	}
	MID, err := db.QueryMemberIDByToken(t)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided"})
		return
	}
	Render(c, gin.H{
		"title":   "Personal Logs",
		"payload": GetAllLogsByMemberID(MID)}, "personal-log-list.html")
}

func GetLog(c *gin.Context) {
	if LogID, err := strconv.Atoi(c.Param("log_id")); err == nil { // Obtain log ID from URL
		t, err := c.Cookie("token")
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		if Token, ok, err := db.QueryTokenExists(t); ok && err == nil && Token == t {
			MID, err := db.QueryMemberIDByToken(t)
			if err != nil {
				log.Println(err)
				c.AbortWithError(http.StatusUnauthorized, err)
			}
			if log, err := GetLogByID(LogID, MID); err == nil {
				Render(c, gin.H{"payload": log}, "personal-log.html")
			} else {
				c.AbortWithError(http.StatusNotFound, err)
			}
		} else {
			// If log is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
		// If invalid log ID, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func ShowLogCreationPage(c *gin.Context) {

}

func CreateLog(c *gin.Context) {

}
