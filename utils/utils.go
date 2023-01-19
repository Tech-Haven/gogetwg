package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tech-haven/gogetwg/configs"
	"github.com/tech-haven/gogetwg/responses"
)

// DESCRIPTION:	returns Wireguard config for the specified clientid

func GetExtClientConf(config *configs.Config, clientid string) ([]byte, *responses.AppError) {
	url := fmt.Sprintf("%s/api/extclients/clients/%s/file", config.NetmakerApiUrl, clientid)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Print(err)
		return nil, &responses.AppError{Error: err, Code: 500, Message: "Error creating HTTP request"}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", configs.MasterKey()))
	resp, err := config.HttpClient.Do(req)
	if err != nil {
		fmt.Print(err)
		return nil, &responses.AppError{Error: err, Code: 500, Message: "Error sending HTTP request"}
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var nmRes responses.NMResponse

		err := json.NewDecoder(resp.Body).Decode(&nmRes)
		if err != nil {
			fmt.Print(err)
			return nil, &responses.AppError{Error: err, Code: 500, Message: "Error decoding response body"}
		}

		if nmRes.Message == "no result found" {
			return nil, &responses.AppError{Error: nil, Code: 404, Message: nmRes.Message}
		}
		return nil, &responses.AppError{Error: nil, Code: nmRes.Code, Message: nmRes.Message}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
		return nil, &responses.AppError{Error: err, Code: 500, Message: "Error decoding response body"}
	}

	return body, nil
}
