package handlers

import "exocomp/engine"
import "fmt"
import "net/http"

func RequestTimeout(session *engine.Session, request *http.Request, response http.ResponseWriter) {

	session.Console.Error(fmt.Sprintf("> %s %s: %d", request.Method, request.URL.Path, http.StatusRequestTimeout))

	content_type, payload := format_error(request, "Request Timeout")

	response.Header().Set("Content-Type", content_type)
	response.WriteHeader(http.StatusRequestTimeout)
	response.Write(payload)

}
