package utils

import "exocomp/schemas"
import "encoding/json"
import "fmt"
import "io"
import "net/http"
import "net/url"
import "slices"

func CheckServer(api_url *url.URL, model_name string) error {

	response, err1 := http.Get(api_url.String())

	if err1 == nil {

		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {

			response_payload, err2 := io.ReadAll(response.Body)

			if err2 == nil {

				schema := schemas.ModelsResponse{}
				err3   := json.Unmarshal(response_payload, &schema)

				if err3 == nil {

					server_type := schema.OwnedBy()

					if server_type == "llamacpp" || server_type == "ollama" || server_type == "vllm" {

						found := false

						for _, model := range schema.Data {

							if model.ID == model_name || slices.Contains(model.Aliases, model_name) {
								found = true
								break
							}

						}

						if found == true {
							return nil
						} else {
							return fmt.Errorf("Error: Model %s not found", model_name)
						}

					} else {
						return fmt.Errorf("Error: Unknown server type")
					}

				} else {
					return fmt.Errorf("Error: %s", err3.Error())
				}

			} else {
				return fmt.Errorf("Error: %s", err2.Error())
			}

		} else {
			return fmt.Errorf("Error: Unexpected HTTP status code %d", response.StatusCode)
		}

	} else {
		return fmt.Errorf("Error: Unexpected HTTP request error %s", err1.Error())
	}

}
