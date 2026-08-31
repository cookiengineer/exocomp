package types

import "exocomp/schemas"
import net_url "net/url"

type Adapter interface {
	Detect(*net_url.URL, string)            bool
	Name()                                  string
	TransformRequest(schemas.ChatRequest)   schemas.ChatRequest
	TransformResponse(schemas.ChatResponse) schemas.ChatResponse
}
