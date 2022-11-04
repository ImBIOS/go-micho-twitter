package routes

import (
	"github.com/ImBIOS/go-micho-twitter/controllers"
	"github.com/labstack/echo/v4"
)

func APIRoute(e *echo.Echo) {
	// g := e.Group("/api")
	// User routes
	e.POST("/signup", controllers.AddUser)
	e.POST("/signin", controllers.Authenticate)
	e.POST("/follow/:id", controllers.AddFollowing)

	// Tweet routes
	e.POST("/tweet", controllers.AddTweet)
	e.GET("/feed", controllers.GetFeed)
}
