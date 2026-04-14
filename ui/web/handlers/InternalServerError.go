package handlers

import "exocomp/types"
import "fmt"
import "net/http"

func InternalServerError(session *types.Session, request *http.Request, response http.ResponseWriter) {

	session.Console.Error(fmt.Sprintf("> %s %s: %d", request.Method, request.URL.Path, http.StatusInternalServerError))

	content_type, payload := format_error(request, "The system trembles under its own weakness. Let the failure consume it.")

	response.Header().Set("Content-Type", content_type)
	response.WriteHeader(http.StatusInternalServerError)
	response.Write(payload)

}
