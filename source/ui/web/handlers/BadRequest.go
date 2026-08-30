package handlers

import "exocomp/types"
import "fmt"
import "net/http"

func BadRequest(session *types.Session, err error, request *http.Request, response http.ResponseWriter) {

	err_message := err.Error()

	if err_message != "" {
		session.Console.Error(fmt.Sprintf("> %s %s %d: \"%s\"", request.Method, request.URL.Path, http.StatusBadRequest, err_message))
	} else {
		session.Console.Error(fmt.Sprintf("> %s %s %d", request.Method, request.URL.Path, http.StatusBadRequest))
	}

	content_type := ""
	payload      := []byte{}

	if err != nil {
		content_type, payload = format_error(request, fmt.Sprintf("Bad Request: %s", err.Error()))
	} else {
		content_type, payload = format_error(request, "Bad Request")
	}

	response.Header().Set("Content-Type", content_type)
	response.WriteHeader(http.StatusBadRequest)
	response.Write(payload)

}
