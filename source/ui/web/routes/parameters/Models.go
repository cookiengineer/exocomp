package parameters

import "exocomp/engine"
import "exocomp/schemas"
import "exocomp/types"
import "exocomp/ui/web/handlers"
import "bytes"
import "encoding/json"
import "fmt"
import "net/http"
import "io"
import "sort"
import "strconv"
import "strings"

var models_cache []schemas.Model

func init() {
	models_cache = make([]schemas.Model, 0)
}

func update_models(config *types.Config) {

	resolved_endpoints := make(map[string]string)
	resolved_models := make(map[string]schemas.Model)

	if config.URL != nil {

		tmp := config.ResolveURL("", "/models")
		resolved_endpoints[tmp.String()] = ""

	}

	for alias, _ := range config.Providers {

		tmp1 := config.ResolveURL(alias, "/models")
		tmp2 := config.ResolveToken(alias)

		resolved_endpoints[tmp1.String()] = tmp2

	}

	for resolved_url, resolved_token := range resolved_endpoints {

		request, err1 := http.NewRequest(
			http.MethodGet,
			resolved_url,
			bytes.NewReader([]byte{}),
		)

		if err1 == nil {

			request.Header.Set("Accept", "application/json")

			if resolved_token != "" {
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resolved_token))
			}

			response, err2 := http.DefaultClient.Do(request)

			if err2 == nil && response.StatusCode == 200 {

				response_payload, err3 := io.ReadAll(response.Body)

				if err3 == nil {

					var schema schemas.ModelsResponse

					err3 := json.Unmarshal(response_payload, &schema)

					if err3 == nil {

						for _, model := range schema.Data {

							model_id := strings.TrimSpace(model.ID)

							if model_id != "" {
								resolved_models[model_id] = model
							}

						}

					}

				}

			}

		}

	}

	if len(resolved_models) > 0 {

		for _, model := range resolved_models {
			models_cache = append(models_cache, model)
		}

	}

}

func Models(session *engine.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodGet {

		if len(models_cache) == 0 {
			update_models(session.Config)
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
