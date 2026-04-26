package types

import "exocomp/schemas"
import utils_chat "exocomp/utils/chat"
import "bytes"
import "encoding/json"
import "fmt"
import "io"
import "net/http"
import "os"
import "strings"
import "sync"

type SessionContext struct {
	Length int `json:"length"`
	Tokens int `json:"tokens"`
}

type Session struct {
	Agent    *Agent             `json:"agent"`
	Config   *Config            `json:"config"`
	Console  *Console           `json:"console"`
	Context  SessionContext     `json:"context"`
	Tools    []*schemas.Tool    `json:"tools"`
	Waiting  bool               `json:"waiting"`
	client   *http.Client       `json:"-"`
	mutex    *sync.RWMutex      `json:"-"`
	tools    map[string]Tool    `json:"-"`
}

func NewSession(agent *Agent, config *Config) *Session {

	session := &Session{
		Agent:    agent,
		Config:   config,
		Console:  NewConsole(os.Stdout, os.Stderr, 0),
		Context:  SessionContext{
			Length: 0,
			Tokens: 0,
		},
		Tools:    make([]*schemas.Tool, 0),
		Waiting:  false,
		client:   &http.Client{},
		mutex:    &sync.RWMutex{},
		tools:    make(map[string]Tool),
	}

	auto_init := false

	session.mutex.Lock()

	if config != nil && config.GetPrompt() != "" {

		session.Agent.Messages = append(session.Agent.Messages, &schemas.Message{
			Role:    "user",
			Content: config.GetPrompt(),
		})

		auto_init = true

	}

	session.Context.Length = session.Config.GetContextLength()

	session.mutex.Unlock()

	if auto_init == true {
		session.Init()
	}

	return session

}

func (session *Session) Init() error {

	if len(session.Agent.Messages) > 0 {
		return session.infer_chat_completions()
	}

	return fmt.Errorf("Session is empty, waiting for LLM system prompt ...")

}

func (session *Session) GetConsoleMessages(from int) []ConsoleMessage {

	if session.Console != nil {
		return session.Console.GetMessages(from)
	} else {
		return []ConsoleMessage{}
	}

}

func (session *Session) GetContextUsage() float64 {

	if session.Context.Length > 0 {
		return float64(float64(session.Context.Tokens) / float64(session.Context.Length)) * 100.0
	} else {
		return 0.0
	}

}

func (session *Session) GetMessages(from int) []*schemas.Message {

    session.mutex.RLock()
    defer session.mutex.RUnlock()

	result := make([]*schemas.Message, 0)

	if len(session.Agent.Messages) > 0 && from < len(session.Agent.Messages) {

		for m := from; m < len(session.Agent.Messages); m++ {
			result = append(result, session.Agent.Messages[m])
		}

	}

	return result

}

func (session *Session) GetTool(identifier string) Tool {

	allowed := false

	for _, tool := range session.Tools {

		if tool.Function.Name == identifier {
			allowed = true
			break
		}

	}

	if allowed == true {

		name     := strings.TrimSpace(identifier[0:strings.Index(identifier, ".")])
		tool, ok := session.tools[name]

		if ok == true {
			return tool
		} else {
			return nil
		}

	} else {
		return nil
	}

}

func (session *Session) SendChatRequest(request schemas.Message) error {

	is_waiting := false

	session.mutex.RLock()
	is_waiting = session.Waiting
	session.mutex.RUnlock()

	if is_waiting == false {

		session.mutex.Lock()
		session.Agent.Messages = append(session.Agent.Messages, &request)
		session.Waiting = true

		session.mutex.Unlock()

		err := session.infer_chat_completions()

		session.mutex.Lock()
		session.Waiting = false
		session.mutex.Unlock()

		if err == nil {
			return nil
		} else {
			return err
		}

	} else {
		return fmt.Errorf("Session is busy, waiting for LLM response ...")
	}

}

func (session *Session) ReceiveChatResponse(response schemas.Message) error {

	if response.Role == "assistant" {

		session.mutex.Lock()
		tmp := &response
		session.Agent.Messages = append(session.Agent.Messages, tmp)
		session.mutex.Unlock()

		if len(response.ToolCalls) > 0 {

			for _, tool_call := range response.ToolCalls {

				identifier, err0 := tool_call.Function.Tool()
				method,     err1 := tool_call.Function.Method()
				arguments,  err2 := tool_call.Function.Arguments()

				if err0 == nil && err1 == nil && err2 == nil {

					tool := session.GetTool(identifier)

					if tool != nil {

						result, err0 := tool.Call(method, arguments)

						if err0 == nil {

							session.mutex.Lock()
							message := &schemas.Message{
								Role:     "tool",
								Content:  strings.TrimSpace(result),
								ToolName: identifier,
							}
							session.Agent.Messages = append(session.Agent.Messages, message)
							session.mutex.Unlock()

						} else {

							session.mutex.Lock()
							message := &schemas.Message{
								Role:     "tool",
								Content:  fmt.Sprintf("Error: %s", strings.TrimSpace(err0.Error())),
								ToolName: identifier,
							}
							session.Agent.Messages = append(session.Agent.Messages, message)
							session.mutex.Unlock()

						}

					} else {

						json_blob, _ := json.MarshalIndent(tool_call, "", "\t")

						session.mutex.Lock()
						message := &schemas.Message{
							Role:     "tool",
							Content:  strings.Join([]string{
								fmt.Sprintf("Error: %s", "Invalid Tool Call"),
								"",
								string(json_blob),
							}, "\n"),
							ToolName: identifier,
						}
						session.Agent.Messages = append(session.Agent.Messages, message)
						session.mutex.Unlock()

					}

				}

			}

			return session.infer_chat_completions()

		} else {
			return nil
		}

	} else {

		session.mutex.Lock()
		tmp := &response
		session.Agent.Messages = append(session.Agent.Messages, tmp)
		session.mutex.Unlock()

		return nil

	}

}

func (session *Session) SetTool(identifier string, tool Tool, schemas []schemas.Tool) {

	if identifier != "" && len(schemas) > 0 && tool != nil {

		session.tools[identifier] = tool

		for _, schema := range schemas {
			session.Tools = append(session.Tools, &schema)
		}

	}

}

func (session *Session) infer_chat_completions() error {

	request_payload, err0 := json.Marshal(schemas.ChatRequest{
		Model:       session.Agent.Model,
		Temperature: session.Agent.Temperature,
		Messages:    session.Agent.Messages,
		Stream:      false,
		Tools:       session.Tools,
		ToolChoice:  "auto",
	})

	if err0 == nil {

		response, err1 := session.client.Post(
			session.Config.ResolveAPI("/v1/chat/completions").String(),
			"application/json",
			bytes.NewReader(request_payload),
		)

		if err1 == nil && response.StatusCode == 200 {

			response_payload, err2 := io.ReadAll(response.Body)

			if err2 == nil {

				var response schemas.ChatResponse

				err3 := json.Unmarshal(response_payload, &response)

				if err3 == nil {

					if response.Usage != nil && response.Usage.PromptTokens != 0 {
						session.Context.Tokens = response.Usage.PromptTokens
					} else {
						session.Context.Tokens = utils_chat.CalculateTokens(session.Agent.Messages)
					}

					if len(response.Choices) > 0 {
						return session.ReceiveChatResponse(response.Choices[0].Message)
					} else {
						return fmt.Errorf("Empty choices, maybe incompatible API?")
					}

				} else {
					return err3
				}

			} else {
				return err2
			}

		} else if err1 == nil && response.StatusCode == 404 {
			return fmt.Errorf("Model %s not found", session.Config.Model)
		} else {
			return err1
		}

	} else {
		return err0
	}

}

