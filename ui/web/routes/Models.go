package routes

import "exocomp/schemas"
import "exocomp/types"
import "exocomp/ui/web/handlers"
import "encoding/json"
import "net/http"
import "io"
import "sort"
import "strconv"
import "strings"

func Models(session *types.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodGet {

		endpoint        := session.Config.ResolvePath("/v1/models")
		response1, err1 := session.Client.Get(endpoint.String())

		if err1 == nil && response1.StatusCode == 200 {

			response1_payload, err2 := io.ReadAll(response1.Body)

			if err2 == nil {

				var schema schemas.ModelsResponse

				err3 := json.Unmarshal(response1_payload, &schema)

				if err3 == nil {

					models := make([]string, 0)

					for _, model := range schema.Data {

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
						handlers.InternalServerError(session, request, response)
					}

				} else {
					handlers.InternalServerError(session, request, response)
				}

			} else {
				handlers.InternalServerError(session, request, response)
			}

		} else if err1 == nil && response1.StatusCode == 404 {
			handlers.NotFound(session, request, response)
		} else {
			handlers.InternalServerError(session, request, response)
		}

	} else {
		handlers.MethodNotAllowed(session, request, response)
	}

}
