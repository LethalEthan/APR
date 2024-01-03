// main.go

package main

import (
	"APR/console"
	"APR/db"
	"APR/routes"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	err := db.DeleteAllTokens()
	if err != nil {
		log.Println(err)
	}
	M, _ := db.QueryMemberInfoByID(3)
	log.Println(M.MemberID, M.FirstName, M.LastName, M.Email, M.Password)
	// Set Gin to production mode
	gin.SetMode(gin.DebugMode)
	router = gin.Default()
	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static/")
	routes.InitialiseRoutes(router)
	// Start serving the application
	go router.Run("127.0.0.1:8080")
	go console.Console()
	//go events.AutoMessage(LDBI, "1092940091275612160")
	console.Shutdown = make(chan os.Signal, 1)
	signal.Notify(console.Shutdown, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-console.Shutdown
	println("Starting Shutdown")
	err = db.DeleteAllTokens()
	if err != nil {
		log.Println(err)
	}
}
