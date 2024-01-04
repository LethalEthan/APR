package routes

import (
	"APR/auth"

	"github.com/gin-gonic/gin"
)

func InitialiseRoutes(router *gin.Engine) {
	router.Use(auth.SetUserStatus())
	router.GET("/", ShowIndexPage) // Handle the index route
	router.GET("/nutrition", ShowNutritionPage)
	router.GET("/pin", ShowPinPage)
	router.GET("/coaching", ShowCoachingPage)
	userRoutes := router.Group("/u") // User group
	{
		userRoutes.GET("/login", auth.EnsureNotLoggedIn(), showLoginPage)
		userRoutes.POST("/login", auth.EnsureNotLoggedIn(), performLogin)
		userRoutes.GET("/logout", auth.EnsureLoggedIn(), logout)
		userRoutes.GET("/register", auth.EnsureNotLoggedIn(), showRegistrationPage)
		userRoutes.POST("/register", auth.EnsureNotLoggedIn(), register)
		userRoutes.GET("/personal-details", auth.EnsureLoggedIn(), showPersonalDetails)
		userRoutes.POST("/personal-details", auth.EnsureLoggedIn(), updatePersonalDetails)
		userRoutes.GET("/profile", auth.EnsureLoggedIn(), showProfile)
	}
	logRoutes := router.Group("/personal-log") // log group
	{
		logRoutes.GET("/view", auth.EnsureLoggedIn(), ShowLogs)
		logRoutes.GET("/view/:log_id", auth.EnsureLoggedIn(), GetLog)
		logRoutes.GET("/create", auth.EnsureLoggedIn(), ShowLogCreationPage)
		logRoutes.POST("/create", auth.EnsureLoggedIn(), CreateLog)
	}
}
