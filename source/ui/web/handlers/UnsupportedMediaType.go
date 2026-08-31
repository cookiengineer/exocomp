package handlers

import "exocomp/engine"
import "fmt"
import "net/http"

func UnsupportedMediaType(session *engine.Session, request *http.Request, response http.ResponseWriter) {

	session.Console.Error(fmt.Sprintf("> %s %s %d", request.Method, request.URL.Path, http.StatusUnsupportedMediaType))

	content_type, payload := format_error(request, "Unsupported Media Type")

	response.Header().Set("Content-Type", content_type)
	response.WriteHeader(http.StatusUnsupportedMediaType)
	response.Write(payload)

}
