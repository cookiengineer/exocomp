package handlers

import "exocomp/engine"
import "fmt"
import "net/http"

func Forbidden(session *engine.Session, request *http.Request, response http.ResponseWriter) {

	session.Console.Error(fmt.Sprintf("> %s %s %d", request.Method, request.URL.Path, http.StatusForbidden))

	content_type, payload := format_error(request, "Forbidden")

	response.Header().Set("Content-Type", content_type)
	response.WriteHeader(http.StatusForbidden)
	response.Write(payload)

}
