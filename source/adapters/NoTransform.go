package adapters

import "exocomp/schemas"
import net_url "net/url"

type NoTransform struct {
}

func (adapter *NoTransform) Name() string {
	return "notransform"
}

func (adapter *NoTransform) Detect(url *net_url.URL, alias string) bool {

	if url != nil {

		if url.Host == "localhost" {
			return true
		}

	}

	return false

}

func (adapter *NoTransform) TransformRequest(request schemas.ChatRequest) schemas.ChatRequest {
	return request
}

func (adapter *NoTransform) TransformResponse(response schemas.ChatResponse) schemas.ChatResponse {
	return response
}
