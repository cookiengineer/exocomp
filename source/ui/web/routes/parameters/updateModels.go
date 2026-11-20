package parameters

import "exocomp/schemas"
import "exocomp/types"
import "bytes"
import "encoding/json"
import "fmt"
import "io"
import "net/http"
import "strings"

func updateModels(config *types.Config) {

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
