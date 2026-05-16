package ollama

import "bytes"
import "encoding/json"
import "io"
import "net/url"
import "net/http"

func GetContextLength(base_url *url.URL, model string) int {

	_, ok := context_lengths[model]

	if ok == false {

		client := &http.Client{}

		endpoint := base_url.ResolveReference(&url.URL{
			Path: "/api/show",
		})

		if endpoint != nil {

			request_body, _ := json.Marshal(ShowRequest{
				Name: model,
			})

			request, err1 := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBuffer(request_body))

			if err1 == nil {

				request.Header.Set("Accept", "application/json")
				request.Header.Set("Content-Type", "application/json")

				response, err2 := client.Do(request)

				if err2 == nil {

					response_payload, err3 := io.ReadAll(response.Body)

					if err3 == nil {

						schema := ShowResponse{}
						err4   := json.Unmarshal(response_payload, &schema)

						if err4 == nil {
							context_lengths[model] = schema.ContextLength()
						}

					}

				}

			}

		}

	}

	return context_lengths[model]

}
