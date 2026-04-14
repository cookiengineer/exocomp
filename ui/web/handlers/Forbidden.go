package handlers

import "exocomp/types"
import "fmt"
import "net/http"

func Forbidden(session *types.Session, request *http.Request, response http.ResponseWriter) {

	session.Console.Error(fmt.Sprintf("> %s %s: %d", request.Method, request.URL.Path, http.StatusForbidden))

	content_type, payload := format_error(request, "The Dark Side forbids this path. You are forbidden from this knowledge.")

	response.Header().Set("Content-Type", content_type)
	response.WriteHeader(http.StatusForbidden)
	response.Write(payload)

}
