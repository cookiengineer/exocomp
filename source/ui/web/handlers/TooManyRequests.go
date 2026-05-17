package handlers

import "exocomp/types"
import "fmt"
import "net/http"

func TooManyRequests(session *types.Session, request *http.Request, response http.ResponseWriter) {

	session.Console.Error(fmt.Sprintf("> %s %s: %d", request.Method, request.URL.Path, http.StatusTooManyRequests))

	content_type, payload := format_error(request, "Too Many Requests")

	response.Header().Set("Content-Type", content_type)
	response.WriteHeader(http.StatusTooManyRequests)
	response.Write(payload)

}
