package this

import (
	"encoding/json"
	"io"
	"main/src/errors"
	services "main/src/service"
	"net/http"
)

func Handler(writer http.ResponseWriter, request *http.Request) {
	sessionToken := request.Header.Get("Authorization")
	isSessionValid := services.ValidateSession(sessionToken)

	if !isSessionValid {
		writer.WriteHeader(401)
		writer.Write(errors.UNAUTHORIZED_BODY)

		return
	}

	if request.Method == http.MethodPost {
		var sendItContent services.TResendBody

		json.NewDecoder(request.Body).Decode(&sendItContent)
		response := services.SendIt(&sendItContent)

		defer response.Body.Close()

		body, _ := io.ReadAll(response.Body)

		writer.WriteHeader(response.StatusCode)
		writer.Write(body)

		return
	}

	writer.WriteHeader(404)
	writer.Write(errors.NOT_FOUND_BODY)
}
