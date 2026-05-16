package ollama

import "bytes"
import "encoding/json"
import "net/http"
import "net/url"

func load_model(base_url *url.URL, model string) bool {

	endpoint := base_url.ResolveReference(&url.URL{
		Path: "/api/generate",
	})

	if endpoint != nil {

		client          := &http.Client{}
		request_body, _ := json.Marshal(GenerateRequest{
			Model:     model,
			Prompt:    "Reply with \"Hello\"!",
			KeepAlive: -1,
		})

		request, err0 := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBuffer(request_body))

		if err0 == nil {

			request.Header.Set("Accept", "application/json")
			request.Header.Set("Content-Type", "application/json")

			response, err1 := client.Do(request)

			if err1 == nil {

				if response.StatusCode == http.StatusOK {
					return true
				}

			}

		}

	}

	return false

}
