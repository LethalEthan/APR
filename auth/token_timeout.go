package auth

import (
	"APR/db"
	"time"
)

// Timeout toklen after 3500 seconds
func TokenTimeout(t string) {
	var timer = time.NewTimer(time.Second * 3500)
	<-timer.C
	db.RevokeToken(t)
}
