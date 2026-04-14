package handlers

import "exocomp/types"
import "fmt"
import "net/http"

func NotFound(session *types.Session, request *http.Request, response http.ResponseWriter) {

	session.Console.Error(fmt.Sprintf("> %s %s: %d", request.Method, request.URL.Path, http.StatusNotFound))

	content_type, payload := format_error(request, "What you seek has vanished into the void. Only shadows remain.")

	response.Header().Set("Content-Type", content_type)
	response.WriteHeader(http.StatusNotFound)
	response.Write(payload)

}
