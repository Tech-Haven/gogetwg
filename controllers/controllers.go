package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tech-haven/gogetwg/configs"
	"github.com/tech-haven/gogetwg/responses"
	"github.com/tech-haven/gogetwg/utils"
)

// ROUTE: 			GET /clients/:clientid, GetClientConfig
//
// DESCRIPTION: Get extclient config of the specified clientID
//
// RESPONSE:		Wireguard config file

func GetExtClientConf(config *configs.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		clientid := c.Param("clientid")

		if clientid == "" {
			return c.JSON(http.StatusBadRequest, responses.HTTPResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "Please specify clientid parameter."}})
		}

		resBody, err := utils.GetExtClientConf(config, clientid)
		if err != nil {
			return c.JSON(err.Code, responses.HTTPResponse{Status: err.Code, Message: "error", Data: &echo.Map{"data": err.Message}})
		}

		return c.JSON(http.StatusOK, responses.HTTPResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": string(resBody)}})
	}
}

// ROUTE: 			POST /clients CreateExtClient
//
// DESCRIPTION: Create a new extclient config. Clientid must be unique.

func CreateExtClient(config *configs.Config) echo.HandlerFunc {
	return func(c echo.Context) error {

		type reqBody struct {
			ClientID string `json:"clientid"`
		}

		jsonBody := new(reqBody)

		if err := c.Bind(&jsonBody); err != nil {
			return c.JSON(http.StatusBadRequest, responses.HTTPResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "Please specify clientid in request body."}})
		}

		url := fmt.Sprintf("%s/api/extclients/clients/%s", config.NetmakerApiUrl, config.NetmakerIngressNodeID)
		values := map[string]string{
			"ClientID": jsonBody.ClientID,
		}
		var response []byte

		body, err := json.Marshal(values)

		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, responses.HTTPResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": "Error creating request body."}})
		}

		_, appErr := utils.NewHTTPRequest(http.MethodPost, url, bytes.NewReader(body), config, response)
		if appErr != nil {
			return c.JSON(appErr.Code, responses.HTTPResponse{Status: appErr.Code, Message: "error", Data: &echo.Map{"data": appErr.Message}})
		}

		resBody, appErr := utils.GetExtClientConf(config, jsonBody.ClientID)
		if appErr != nil {
			return c.JSON(appErr.Code, responses.HTTPResponse{Status: appErr.Code, Message: "error", Data: &echo.Map{"data": appErr.Message}})
		}

		return c.JSON(http.StatusOK, responses.HTTPResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": string(resBody)}})
	}
}

// ROUTE: 			GET /ping, Ping
//
// DESCRIPTION: Healthcheck
//
// RESPONSE:		String

func Ping() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, responses.HTTPResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": "Pong!"}})
	}
}
