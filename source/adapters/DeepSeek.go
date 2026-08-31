package adapters

import "exocomp/schemas"
import net_url "net/url"

type DeepSeek struct {
}

func (adapter *DeepSeek) Name() string {
	return "deepseek"
}

func (adapter *DeepSeek) Detect(url *net_url.URL, alias string) bool {

	if url != nil {

		if url.Host == "api.deepseek.com" {
			return true
		}

	}

	available_models := []string{
		"deepseek-v4-flash",
		"deepseek-v4-flash-vision",
		"deepseek-v4-pro",
	}

	for _, model := range available_models {

		if model == alias {
			return true
		}

	}

	return false

}

func (adapter *DeepSeek) TransformRequest(request schemas.ChatRequest) schemas.ChatRequest {

	// TODO: Modify tools to be what_Ever instead of what.Ever

	return request

}

func (adapter *DeepSeek) TransformResponse(response schemas.ChatResponse) schemas.ChatResponse {

	// TODO: Modify tools to be what.Ever instead of what_Ever

	return response

}
