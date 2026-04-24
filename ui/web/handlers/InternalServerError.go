package handlers

import "exocomp/types"
import "fmt"
import "net/http"

func InternalServerError(session *types.Session, err error, request *http.Request, response http.ResponseWriter) {

	session.Console.Error(fmt.Sprintf("> %s %s: %d", request.Method, request.URL.Path, http.StatusInternalServerError))

	content_type := ""
	payload      := []byte{}

	if err != nil {
		content_type, payload = format_error(request, fmt.Sprintf("Internal Server Error: %s", err.Error()))
	} else {
		content_type, payload = format_error(request, "Internal Server Error")
	}

	response.Header().Set("Content-Type", content_type)
	response.WriteHeader(http.StatusInternalServerError)
	response.Write(payload)

}
