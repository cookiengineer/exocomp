package handlers

import "exocomp/types"
import "fmt"
import "net/http"

func Unauthorized(session *types.Session, request *http.Request, response http.ResponseWriter) {

	session.Console.Error(fmt.Sprintf("> %s %s: %d", request.Method, request.URL.Path, http.StatusUnauthorized))

	content_type, payload := format_error(request, "The Dark Side does not recognize you. You lack the power to authenticate.")

	response.Header().Set("Content-Type", content_type)
	response.WriteHeader(http.StatusUnauthorized)
	response.Write(payload)

}
