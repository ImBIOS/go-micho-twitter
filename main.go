package main

import (
	"github.com/ImBIOS/go-micho-twitter/configs"
	"github.com/ImBIOS/go-micho-twitter/controllers"
	"github.com/ImBIOS/go-micho-twitter/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	// Init Echo
	e := echo.New()

	// Set log level to error
	e.Logger.SetLevel(log.ERROR)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(controllers.Key),
		Skipper: func(c echo.Context) bool {
			// Skip authentication for signup and signin requests
			return c.Path() == "/signin" || c.Path() == "/signup"
		},
	}))

	// Custom JWT Error
	middleware.ErrJWTMissing.Code = 401
	middleware.ErrJWTMissing.Message = "Unauthorized"

	// Database connection
	configs.ConnectDB()

	// Routes
	routes.APIRoute(e)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
