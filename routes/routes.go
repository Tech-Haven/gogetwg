package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/tech-haven/gogetwg/configs"
	"github.com/tech-haven/gogetwg/controllers"
)

func Routes(e *echo.Echo, configuration *configs.Config) {
	e.GET("/ping", controllers.Ping())
	e.GET("/clients/:clientid", controllers.GetExtClientConf(configuration))
	e.POST("/clients", controllers.CreateExtClient(configuration))
}
