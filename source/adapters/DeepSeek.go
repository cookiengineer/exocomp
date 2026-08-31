package adapters

import "exocomp/schemas"
import net_url "net/url"
import "strings"

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

	for m, message := range request.Messages {

		message.ToolName = strings.ReplaceAll(message.ToolName, ".", "_")

		for t, toolcall := range message.ToolCalls {

			toolcall.Function.Name = strings.ReplaceAll(toolcall.Function.Name, ".", "_")

			message.ToolCalls[t] = toolcall

		}

		request.Messages[m] = message

	}

	for t, tool := range request.Tools {

		tool.Function.Name = strings.ReplaceAll(tool.Function.Name, ".", "_")

		request.Tools[t] = tool

	}

	return request

}

func (adapter *DeepSeek) TransformResponse(response schemas.ChatResponse) schemas.ChatResponse {

	for c, choice := range response.Choices {

		choice.Message.ToolName = strings.ReplaceAll(choice.Message.ToolName, "_", ".")

		for t, toolcall := range choice.Message.ToolCalls {

			toolcall.Function.Name = strings.ReplaceAll(toolcall.Function.Name, "_", ".")

			choice.Message.ToolCalls[t] = toolcall

		}

		response.Choices[c] = choice

	}

	return response

}
