package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context, data gin.H, templateName string) {
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"]) // Respond with JSON
	case "application/xml":
		c.XML(http.StatusOK, data["payload"]) // Respond with XML
	default:
		c.HTML(http.StatusOK, templateName, data) // Respond with HTML
	}
}
