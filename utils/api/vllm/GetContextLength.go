package vllm

import "encoding/json"
import "io"
import "net/url"
import "net/http"

func GetContextLength(base_url *url.URL, model string) int {

	_, ok := context_lengths[model]

	if ok == false {

		client := &http.Client{}

		endpoint := base_url.ResolveReference(&url.URL{
			Path: "/get_model_config",
		})

		if endpoint != nil {

			request, err1 := http.NewRequest(http.MethodGet, endpoint.String(), nil)

			if err1 == nil {

				request.Header.Set("Accept", "application/json")

				response, err2 := client.Do(request)

				if err2 == nil {

					response_payload, err3 := io.ReadAll(response.Body)

					if err3 == nil {

						schema := ModelConfigResponse{}
						err4   := json.Unmarshal(response_payload, &schema)

						if err4 == nil {
							context_lengths[model] = schema.ContextLength(model)
						}

					}

				}

			}

		}

	}

	return context_lengths[model]

}
