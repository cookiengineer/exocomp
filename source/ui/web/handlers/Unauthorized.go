package handlers

import "exocomp/types"
import "fmt"
import "net/http"

func Unauthorized(session *types.Session, err error, request *http.Request, response http.ResponseWriter) {

	err_message := err.Error()

	if err_message != "" {
		session.Console.Error(fmt.Sprintf("> %s %s %d: \"%s\"", request.Method, request.URL.Path, http.StatusUnauthorized, err_message))
	} else {
		session.Console.Error(fmt.Sprintf("> %s %s %d", request.Method, request.URL.Path, http.StatusUnauthorized))
	}

	content_type := ""
	payload      := []byte{}

	if err != nil {
		content_type, payload = format_error(request, fmt.Sprintf("Unauthorized : %s", err.Error()))
	} else {
		content_type, payload = format_error(request, "Unauthorized")
	}

	response.Header().Set("Content-Type", content_type)
	response.WriteHeader(http.StatusUnauthorized)
	response.Write(payload)

}
