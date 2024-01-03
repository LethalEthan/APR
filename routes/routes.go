package routes

import (
	"APR/auth"

	"github.com/gin-gonic/gin"
)

func InitialiseRoutes(router *gin.Engine) {
	router.Use(auth.SetUserStatus())
	router.GET("/", ShowIndexPage)   // Handle the index route
	userRoutes := router.Group("/u") // User group
	{
		userRoutes.GET("/login", auth.EnsureNotLoggedIn(), showLoginPage)
		userRoutes.POST("/login", auth.EnsureNotLoggedIn(), performLogin)
		userRoutes.GET("/logout", auth.EnsureLoggedIn(), logout)
		userRoutes.GET("/register", auth.EnsureNotLoggedIn(), showRegistrationPage)
		userRoutes.POST("/register", auth.EnsureNotLoggedIn(), register)
		router.GET("/personal-details", auth.EnsureLoggedIn(), showPersonalDetails)
	}
}
