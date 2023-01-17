package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tech-haven/gogetwg/configs"
	"github.com/tech-haven/gogetwg/routes"
)

func main() {
	e := echo.New()

	// Load .env file
	configs.LoadEnv()

	// Get configurations
	configuration, err := configs.New()
	if err != nil {
		log.Fatal("Error creating configuration", err)
	}

	// Routes
	routes.Routes(e, configuration)

	// Token auth middleware
	e.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		return key == configs.AuthSecret(), nil
	}))

	// Start server
	e.Logger.Fatal(e.Start(":1324"))
}
