package handlers

import "exocomp/types"
import "fmt"
import "net/http"

func MethodNotAllowed(session *types.Session, request *http.Request, response http.ResponseWriter) {

	session.Console.Error(fmt.Sprintf("> %s %s: %d", request.Method, request.URL.Path, http.StatusMethodNotAllowed))

	content_type, payload := format_error(request, "Method Not Allowed")

	response.Header().Set("Content-Type", content_type)
	response.WriteHeader(http.StatusMethodNotAllowed)
	response.Write(payload)

}
