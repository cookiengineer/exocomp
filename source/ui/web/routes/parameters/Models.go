package parameters

import "exocomp/engine"
import "exocomp/schemas"
import "exocomp/ui/web/handlers"
import "encoding/json"
import "net/http"
import "sort"
import "strconv"
import "strings"

var models_cache []schemas.Model

func init() {
	models_cache = make([]schemas.Model, 0)
}

func Models(session *engine.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodGet {

		if len(models_cache) == 0 {
			updateModels(session.Config)
		}

		models := make([]string, 0)

		for _, model := range models_cache {

			model_id := strings.TrimSpace(model.ID)

			if model_id != "" {
				models = append(models, model_id)
			}

		}

		sort.Strings(models)

		response_payload, err4 := json.MarshalIndent(models, "", "\t")

		if err4 == nil {

			response.Header().Set("Content-Type", "application/json")
			response.Header().Set("Content-Length", strconv.Itoa(len(response_payload)))
			response.WriteHeader(http.StatusOK)
			response.Write(response_payload)

		} else {
			handlers.InternalServerError(session, err4, request, response)
		}

	} else {
		handlers.MethodNotAllowed(session, request, response)
	}

}
