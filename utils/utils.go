package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tech-haven/gogetwg/configs"
	"github.com/tech-haven/gogetwg/responses"
)

// DESCRIPTION:	returns Wireguard config for the specified clientid

func GetExtClientConf(config *configs.Config, clientid string) ([]byte, *responses.AppError) {
	url := fmt.Sprintf("%s/api/extclients/clients/%s/file", config.NetmakerApiUrl, clientid)
	var response []byte

	res, appErr := NewHTTPRequest("GET", url, nil, config, response)
	if appErr != nil {
		return nil, appErr
	}

	return res, nil
}

// DESCRIPTION:	returns http response
func NewHTTPRequest[T []byte](method string, url string, body io.Reader, config *configs.Config, responseType T) (T, *responses.AppError) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Println(err)
		return responseType, &responses.AppError{Error: err, Code: 500, Message: "Error creating HTTP request"}
	}

	// set Netmaker Auth header for all requests
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", configs.MasterKey()))

	// set content type for POST requests
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	// log the request
	log.Printf("%s %s\n", method, req.URL.String())

	res, err := config.HttpClient.Do(req)
	if err != nil {
		log.Println(err)
		return responseType, &responses.AppError{Error: err, Code: 500, Message: "Error sending HTTP request"}
	}

	defer res.Body.Close()

	// handle Netmaker error
	if res.StatusCode != http.StatusOK {
		var nmRes responses.NetmakerErrorResponse

		err := json.NewDecoder(res.Body).Decode(&nmRes)
		if err != nil {
			log.Println(err)
			return responseType, &responses.AppError{Error: err, Code: 500, Message: "Error decoding response body"}
		}

		if nmRes.Message == "no result found" {
			return responseType, &responses.AppError{Error: nil, Code: 404, Message: nmRes.Message}
		}
		return responseType, &responses.AppError{Error: nil, Code: nmRes.Code, Message: nmRes.Message}
	}

	// Netmaker returns raw text for config requests
	if res.Header.Get("Content-Type") != "application/json" {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			return responseType, &responses.AppError{Error: err, Code: 500, Message: "Error reading response body"}
		}

		return body, nil
	}

	var responseObject T
	err = json.NewDecoder(res.Body).Decode(&responseObject)

	// Netmaker sometimes returns an empty body
	if err == io.EOF {
		return responseObject, nil
	}

	if err != nil {
		log.Println(err)
		return responseType, &responses.AppError{Error: err, Code: 500, Message: err.Error()}
	}

	return responseObject, nil
}
