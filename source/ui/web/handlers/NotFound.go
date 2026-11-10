package handlers

import "exocomp/engine"
import "fmt"
import "net/http"

func NotFound(session *engine.Session, request *http.Request, response http.ResponseWriter) {

	session.Console.Error(fmt.Sprintf("> %s %s %d", request.Method, request.URL.Path, http.StatusNotFound))

	content_type, payload := format_error(request, "Not Found")

	response.Header().Set("Content-Type", content_type)
	response.WriteHeader(http.StatusNotFound)
	response.Write(payload)

}
