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

type JsonBody struct {
	ClientID string `json:"ClientID"`
}

// ROUTE: 			POST /clients CreateExtClient
//
// DESCRIPTION: Create a new extclient config. Clientid must be unique.

func CreateExtClient(config *configs.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		jsonBody := new(JsonBody)

		if err := c.Bind(&jsonBody); err != nil {
			return c.JSON(http.StatusBadRequest, responses.HTTPResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "Please specify ClientID."}})
		}

		url := fmt.Sprintf("%s/api/extclients/clients/%s", config.NetmakerApiUrl, config.NetmakerIngressNodeID)
		values := map[string]string{
			"ClientID": jsonBody.ClientID,
		}
		reqBody, err := json.Marshal(values)

		if err != nil {
			log.Fatalf("Failed to marshal request body: %s", err)
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		if err != nil {
			fmt.Print(err)
			return c.JSON(http.StatusBadRequest, responses.HTTPResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "Something went wrong."}})
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", configs.MasterKey()))
		req.Header.Set("Content-Type", "application/json")
		resp, err := config.HttpClient.Do(req)
		if err != nil {
			fmt.Print(err)
			return c.JSON(http.StatusBadRequest, responses.HTTPResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "Something went wrong."}})
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			var nmRes responses.NMResponse

			err = json.NewDecoder(resp.Body).Decode(&nmRes)
			if err != nil {
				log.Fatalln(err)
			}
			return c.JSON(http.StatusOK, responses.HTTPResponse{Status: http.StatusOK, Message: "error", Data: &echo.Map{"data": nmRes}})
		}

		resBody, e := utils.GetExtClientConf(config, jsonBody.ClientID)
		if e != nil {
			return c.JSON(e.Code, responses.HTTPResponse{Status: e.Code, Message: "error", Data: &echo.Map{"data": e.Message}})
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
