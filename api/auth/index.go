package auth

import (
	"encoding/json"
	"main/src/errors"
	auth "main/src/service"
	"net/http"
	"strings"
)

type TAuthRequest struct {
	ApiKey string `json:"api-key"`
}

type TAuthResponse struct {
	Session string `json:"session"`
}

func Handler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	writer.Header().Set("Content-type", "application/json")

	var authRequest TAuthRequest

	err := json.NewDecoder(request.Body).Decode(&authRequest)

	if err != nil {
		writer.WriteHeader(400)
		writer.Write(errors.BAD_REQUEST_BODY)

		return
	}

	if request.Method == http.MethodDelete {
		sessionString := request.Header.Get("Authorization")
		sessionString = strings.Replace(sessionString, "Bearer ", "", 1)

		auth.RemoveSession(sessionString)

		return
	}

	if request.Method == http.MethodPost {
		sessionToken := auth.CreateSession(authRequest.ApiKey)
		sessionResponse := TAuthResponse{
			Session: sessionToken,
		}

		json.NewEncoder(writer).Encode(&sessionResponse)

		return
	}

	writer.WriteHeader(404)
	writer.Write(errors.NOT_FOUND_BODY)
}
