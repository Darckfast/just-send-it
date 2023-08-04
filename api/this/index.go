package this

import (
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

	writer.WriteHeader(404)
	writer.Write(errors.NOT_FOUND_BODY)
}
