package types

import "exocomp/agents"
import "exocomp/schemas"
import "exocomp/tools"
import "bytes"
import "encoding/json"
import "fmt"
import "io"
import "net/http"
import "os"
import "strings"
import "sync"

type Session struct {
	Agent    *agents.Agent
	Config   *Config
	Console  *Console
	Client   *http.Client
	Messages []*schemas.Message
	Models   []*schemas.Model
	Tools    []*schemas.Tool
	Waiting  bool
	mutex    *sync.RWMutex
	tools    map[string]tools.Tool
}

func NewSession(agent *agents.Agent, config *Config) *Session {

	session := &Session{
		Agent:    agent,
		Config:   config,
		Console:  NewConsole(os.Stdout, os.Stderr, 0),
		Client:   &http.Client{},
		Messages: make([]*schemas.Message, 0),
		Tools:    make([]*schemas.Tool, 0),
		Waiting:  false,
		mutex:    &sync.RWMutex{},
		tools:    make(map[string]tools.Tool),
	}

	if len(agent.Tools) > 0 {
		session.Tools = tools.EncodeSchema(agent.Tools)
	}

	auto_init := false

	session.mutex.Lock()

	if agent != nil {

		system_message := &schemas.Message{
			Role:    "system",
			Content: agent.GetPrompt(),
		}

		session.Messages = append(session.Messages, system_message)

	}

	if config != nil && config.GetPrompt() != "" {

		user_message := &schemas.Message{
			Role:    "user",
			Content: config.GetPrompt(),
		}

		session.Messages = append(session.Messages, user_message)
		auto_init = true

	}

	session.mutex.Unlock()

	if auto_init == true {
		session.Init()
	}

	return session

}

func (session *Session) Init() error {

	if len(session.Messages) > 0 {
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

func (session *Session) GetMessages(from int) []*schemas.Message {

    session.mutex.RLock()
    defer session.mutex.RUnlock()

	result := make([]*schemas.Message, 0)

	if len(session.Messages) > 0 && from < len(session.Messages) {

		for m := from; m < len(session.Messages); m++ {
			result = append(result, session.Messages[m])
		}

	}

	return result

}

func (session *Session) GetTool(identifier string) tools.Tool {

	allowed := false

	for _, tool := range session.Tools {

		if tool.Function.Name == identifier {
			allowed = true
			break
		}

	}

	if allowed == true {

		name  := strings.TrimSpace(identifier[0:strings.Index(identifier, ".")])
		_, ok := session.tools[name]

		if ok == false {


			switch name {
			case "agents":
				session.tools[name] = tools.Tool(tools.NewAgents(session.Config.Agent, session.Config.Sandbox, session.Config.Playground))
			case "bugs":
				session.tools[name] = tools.Tool(tools.NewBugs(session.Config.Agent, session.Config.Sandbox, session.Config.Playground))
			case "changelog":
				session.tools[name] = tools.Tool(tools.NewChangelog(session.Config.Agent, session.Config.Sandbox, session.Config.Playground))
			case "files":
				session.tools[name] = tools.Tool(tools.NewFiles(session.Config.Agent, session.Config.Sandbox))
			case "programs":
				session.tools[name] = tools.Tool(tools.NewPrograms(session.Config.Agent, session.Config.Sandbox, session.Agent.Programs))
			case "requirements":
				session.tools[name] = tools.Tool(tools.NewRequirements(session.Config.Agent, session.Config.Sandbox, session.Config.Playground))
			default:
				session.tools[name] = nil
			}

		}

		return session.tools[name]

	} else {
		return nil
	}

}

func (session *Session) SendChatRequest(raw schemas.Message) error {

	is_waiting := false

	session.mutex.RLock()
	is_waiting = session.Waiting
	session.mutex.RUnlock()

	if is_waiting == false {

		session.mutex.Lock()

		message := &raw
		session.Messages = append(session.Messages, message)
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
		msg := &response
		session.Messages = append(session.Messages, msg)
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
							session.Messages = append(session.Messages, message)
							session.mutex.Unlock()

						} else {

							session.mutex.Lock()
							message := &schemas.Message{
								Role:     "tool",
								Content:  fmt.Sprintf("Error: %s", strings.TrimSpace(err0.Error())),
								ToolName: identifier,
							}
							session.Messages = append(session.Messages, message)
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
						session.Messages = append(session.Messages, message)
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
		msg := &response
		session.Messages = append(session.Messages, msg)
		session.mutex.Unlock()

		return nil

	}

}

func (session *Session) infer_chat_completions() error {

	request_payload, err0 := json.Marshal(schemas.ChatRequest{
		Model:       session.Agent.Model,
		Temperature: session.Agent.Temperature,
		Messages:    session.Messages,
		Stream:      false,
		Tools:       session.Tools,
	})

	if err0 == nil {

		endpoint := session.Config.ResolvePath("/v1/chat/completions")

		response, err1 := session.Client.Post(
			endpoint.String(),
			"application/json",
			bytes.NewReader(request_payload),
		)

		if err1 == nil && response.StatusCode == 200 {

			response_payload, err2 := io.ReadAll(response.Body)

			if err2 == nil {

				var response schemas.ChatResponse

				err3 := json.Unmarshal(response_payload, &response)

				if err3 == nil {

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

