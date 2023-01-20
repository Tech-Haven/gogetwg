package responses

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTTPResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    *echo.Map `json:"data"`
}

type AppError struct {
	Error   error
	Code    int
	Message string
}

type NMResponse struct {
	Code    int
	Message string
}

func (r NMResponse) HandleError(resp *http.Response) *AppError {
	if resp.StatusCode != 200 {
		err := json.NewDecoder(resp.Body).Decode(&r)
		if err != nil {
			log.Println(err)
			return &AppError{Error: err, Code: 500, Message: "Error decoding response body"}
		}

		if r.Message == "no result found" {
			return &AppError{Error: nil, Code: 404, Message: r.Message}
		}
		return &AppError{Error: nil, Code: r.Code, Message: r.Message}
	}
	return nil
}
