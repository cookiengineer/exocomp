package handlers

import "exocomp/types"
import "fmt"
import "net/http"

func RequestTimeout(session *types.Session, request *http.Request, response http.ResponseWriter) {

	session.Console.Error(fmt.Sprintf("> %s %s: %d", request.Method, request.URL.Path, http.StatusRequestTimeout))

	content_type, payload := format_error(request, "The Dark Side grows impatient. Your hesitation has cost you.")

	response.Header().Set("Content-Type", content_type)
	response.WriteHeader(http.StatusRequestTimeout)
	response.Write(payload)

}
