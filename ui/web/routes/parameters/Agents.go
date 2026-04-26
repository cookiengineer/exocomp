package parameters

import "exocomp/agents"
import "exocomp/types"
import "exocomp/ui/web/handlers"
import "encoding/json"
import "net/http"
import "sort"
import "strconv"

func Agents(session *types.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodGet {

		types := make([]string, 0)

		for _, typ := range agents.Types {
			types = append(types, typ)
		}

		sort.Strings(types)

		response_payload, err1 := json.MarshalIndent(types, "", "\t")

		if err1 == nil {

			response.Header().Set("Content-Type", "application/json")
			response.Header().Set("Content-Length", strconv.Itoa(len(response_payload)))
			response.WriteHeader(http.StatusOK)
			response.Write(response_payload)

		} else {
			handlers.InternalServerError(session, err1, request, response)
		}

	} else {
		handlers.MethodNotAllowed(session, request, response)
	}

}
