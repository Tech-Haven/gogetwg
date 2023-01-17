package responses

import "github.com/labstack/echo/v4"

type HTTPResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    *echo.Map `json:"data"`
}

type NMResponse struct {
	Code    int
	Message string
}
