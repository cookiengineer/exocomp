package handlers

import "exocomp/types"
import "net/http"

func SeeOther(session *types.Session, request *http.Request, response http.ResponseWriter, location string) {

	response.Header().Set("Location", location)
	response.WriteHeader(http.StatusSeeOther)
	response.Write([]byte("I hate this place... and you, Mr. Anderson."))

}

