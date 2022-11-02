package routes

import (
	"github.com/ImBIOS/go-micho-twitter/controllers"
	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	e.POST("/signup", controllers.Signup)
	e.POST("/signin", controllers.Signin)
}
