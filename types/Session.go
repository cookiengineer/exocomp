package types

import "exocomp/schemas"
import utils_chat "exocomp/utils/chat"
import "bytes"
import "encoding/json"
import "fmt"
import "io"
import "net/http"
import "os"
import "sort"
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

	session.mutex.Lock()

	if config != nil && config.GetPrompt() != "" {

		session.Agent.Messages = append(session.Agent.Messages, &schemas.Message{
			Role:    "user",
			Content: config.GetPrompt(),
		})

	}

	session.Context.Length = session.Config.GetContextLength()

	session.mutex.Unlock()

	return session

}

func (session *Session) Init() error {

	// NOTE: First Message is System Prompt
	if len(session.Agent.Messages) > 0 {
		return session.infer_chat_completions()
	}

	return fmt.Errorf("Session is empty, waiting for LLM system prompt ...")

}

func (session *Session) CallTool(name string, method string, arguments map[string]any) error {

	tool := session.GetTool(name)

	if tool != nil {

		result, err0 := tool.Call(method, arguments)

		if name == "skills" && method == "Load" {

			tool_message := ""

			if err0 == nil {

				skill_name,    ok1  := arguments["name"].(string)
				skill_content, err1 := tool.Get(skill_name)
				skill,         ok2  := skill_content.(*Skill)

				if ok1 == true && err1 == nil && ok2 == true {

					err2 := session.LoadSkill(skill_name, skill)

					if err2 == nil {
						tool_message = strings.TrimSpace(result)
					} else {
						tool_message = fmt.Sprintf("Error: skills.Load: %s", err2.Error())
					}

				} else {
					tool_message = fmt.Sprintf("Error: skills.Load: %s", "Attempt to escape policies")
				}

			} else {
				tool_message = fmt.Sprintf("Error: skills.Load: %s", strings.TrimSpace(err0.Error()))
			}

			session.mutex.Lock()
			message := &schemas.Message{
				Role:     "tool",
				Content:  tool_message,
				ToolName: name,
			}
			session.Agent.Messages = append(session.Agent.Messages, message)
			session.mutex.Unlock()

			if strings.HasPrefix(tool_message, "Error:") {
				return fmt.Errorf(tool_message)
			} else {
				return nil
			}

		} else if name == "skills" && method == "Unload" {

			tool_message := ""

			if err0 == nil {

				skill_name,    ok1  := arguments["name"].(string)
				skill_content, err1 := tool.Get(skill_name)
				skill,         ok2  := skill_content.(*Skill)

				if ok1 == true && err1 == nil && ok2 == true {

					err2 := session.UnloadSkill(skill_name, skill)

					if err2 == nil {
						tool_message = strings.TrimSpace(result)
					} else {
						tool_message = fmt.Sprintf("Error: skills.Unload: %s", err2.Error())
					}

				} else {
					tool_message = fmt.Sprintf("Error: skills.Unload: %s", "Attempt to escape policies")
				}

			} else {
				tool_message = fmt.Sprintf("Error: skills.Unload: %s", strings.TrimSpace(err0.Error()))
			}

			session.mutex.Lock()
			message := &schemas.Message{
				Role:     "tool",
				Content:  tool_message,
				ToolName: name,
			}
			session.Agent.Messages = append(session.Agent.Messages, message)
			session.mutex.Unlock()

			if strings.HasPrefix(tool_message, "Error:") {
				return fmt.Errorf(tool_message)
			} else {
				return nil
			}

		} else {

			tool_message := ""

			if err0 == nil {
				tool_message = strings.TrimSpace(result)
			} else {
				tool_message = fmt.Sprintf("Error: %s", strings.TrimSpace(err0.Error()))
			}

			session.mutex.Lock()
			message := &schemas.Message{
				Role:     "tool",
				Content:  tool_message,
				ToolName: name,
			}
			session.Agent.Messages = append(session.Agent.Messages, message)
			session.mutex.Unlock()

			if strings.HasPrefix(tool_message, "Error:") {
				return fmt.Errorf(tool_message)
			} else {
				return nil
			}

		}

	} else {

		args_blob, _ := json.Marshal(arguments)
		json_blob, _ := json.Marshal(schemas.ToolCall{
			Type:     "function",
			Function: schemas.ToolCallFunction{
				Name:         name,
				ArgumentsRaw: args_blob,
			},
		})

		session.mutex.Lock()
		message := &schemas.Message{
			Role:     "tool",
			Content:  strings.Join([]string{
				fmt.Sprintf("Error: Tool \"%s\" doesn't exist.", name),
				"",
				string(json_blob),
			}, "\n"),
			ToolName: name,
		}
		session.Agent.Messages = append(session.Agent.Messages, message)
		session.mutex.Unlock()

		return fmt.Errorf("Error: Tool \"%s\" doesn't exist.", name)

	}

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

func (session *Session) GetTool(name string) Tool {

	allowed := false

	for _, tool := range session.Tools {

		if tool.Function.Name == name {
			allowed = true
			break
		}

	}

	if allowed == true {

		namespace := strings.TrimSpace(name[0:strings.Index(name, ".")])
		tool, ok  := session.tools[namespace]

		if ok == true {
			return tool
		} else {
			return nil
		}

	} else {
		return nil
	}

}

func (session *Session) GetToolNames() []string {

	result := make([]string, 0)

	for _, tool := range session.Tools {
		result = append(result, tool.Function.Name)
	}

	sort.Strings(result)

	return result

}

func (session *Session) GetToolSchema(name string) *schemas.Tool {

	var found *schemas.Tool = nil

	for _, tool := range session.Tools {

		if tool.Function.Name == name {
			found = tool
			break
		}

	}

	return found

}

func (session *Session) LoadSkill(name string, skill *Skill) error {

	index            := int(-1)
	missing_programs := make([]string, 0)
	missing_tools    := make([]string, 0)

	session.mutex.Lock()

	for m, message := range session.Agent.Messages {

		if message.Role == "system" && message.Content == skill.Body {
			index = m
			break
		}

	}

	session.mutex.Unlock()

	if len(skill.AllowedPrograms) > 0 {

		for _, program_name := range skill.AllowedPrograms {

			found := false

			for _, program := range session.Agent.AllowedPrograms {

				if program == program_name {
					found = true
					break
				}

			}

			if found == false {
				missing_programs = append(missing_programs, program_name)
			}

		}

	}
	if len(skill.AllowedTools) > 0 {

		for _, tool_name := range skill.AllowedTools {

			found := false

			for _, tool := range session.Agent.AllowedTools {

				if tool == tool_name {
					found = true
					break
				}

			}

			if found == false {
				missing_tools = append(missing_tools, tool_name)
			}

		}

	}

	if index == -1 {

		if len(missing_tools) == 0 {

			system_messages := make([]*schemas.Message, 0)
			other_messages  := make([]*schemas.Message, 0)

			session.mutex.Lock()

			for _, message := range session.Agent.Messages {

				if message.Role == "system" {
					system_messages = append(system_messages, message)
				} else {
					other_messages = append(other_messages, message)
				}

			}

			system_messages = append(system_messages, &schemas.Message{
				Role:    "system",
				Content: skill.Body,
			})
			session.Agent.Messages = append(system_messages, other_messages...)

			session.mutex.Unlock()

			return nil

		} else {
			return fmt.Errorf("Session.LoadSkill: Can't load Skill because of missing Tools %s", strings.Join(missing_tools, " and "))
		}

	} else {
		return fmt.Errorf("Session.LoadSkill: %s", "Skill is already loaded.")
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

				name,      err0 := tool_call.Function.ToName()
				method,    err1 := tool_call.Function.ToMethod()
				arguments, err2 := tool_call.Function.ToArguments()

				if err0 == nil && err1 == nil && err2 == nil {
					session.CallTool(name, method, arguments)
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

func (session *Session) SetTool(identifier string, tool Tool, schemas []schemas.Tool) {

	if identifier != "" && len(schemas) > 0 && tool != nil {

		session.tools[identifier] = tool

		for _, schema := range schemas {
			session.Tools = append(session.Tools, &schema)
		}

	}

}

func (session *Session) UnloadSkill(name string, skill *Skill) error {

	index := int(-1)

	session.mutex.Lock()

	for m, message := range session.Agent.Messages {

		if message.Role == "system" && message.Content == skill.Body {
			index = m
			break
		}

	}

	session.mutex.Unlock()

	if index != -1 {

		system_messages := make([]*schemas.Message, 0)
		other_messages  := make([]*schemas.Message, 0)

		session.mutex.Lock()

		for _, message := range session.Agent.Messages {

			if message.Role == "system" {

				if message.Content != skill.Body {
					system_messages = append(system_messages, message)
				}

			} else {
				other_messages = append(other_messages, message)
			}

		}

		session.Agent.Messages = append(system_messages, other_messages...)

		session.mutex.Unlock()

		return nil

	} else {
		return fmt.Errorf("Session.UnloadSkill: %s", "Skill is already unloaded.")
	}

}

func (session *Session) infer_chat_completions() error {

	request_payload, err0 := json.MarshalIndent(schemas.ChatRequest{
		Model:       session.Agent.Model,
		Temperature: session.Agent.Temperature,
		Messages:    session.Agent.Messages,
		Stream:      false,
		Tools:       session.Tools,
		ToolChoice:  "auto",
		Options:     nil,
		// Options:     &schemas.Options{
		// 	NumContext: 262144,
		// 	NumPredict: 8192,
		// },
	}, "", "\t")

	if err0 == nil {

		response, err1 := session.client.Post(
			session.Config.ResolveAPI("/v1/chat/completions").String(),
			"application/json",
			bytes.NewReader(request_payload),
		)

		if session.Config.Debug == true {

			var tmp1 interface{}
			json.Unmarshal(request_payload, &tmp1)
			tmp2, _ := json.MarshalIndent(tmp1, "", "\t")

			os.WriteFile(session.Config.Sandbox + "/.exocomp-request.json", tmp2, 0666)

		}

		if err1 == nil && response.StatusCode == 200 {

			response_payload, err2 := io.ReadAll(response.Body)

			if session.Config.Debug == true {

				var tmp1 interface{}
				json.Unmarshal(response_payload, &tmp1)
				tmp2, _ := json.MarshalIndent(tmp1, "", "\t")

				os.WriteFile(session.Config.Sandbox + "/.exocomp-response.json", tmp2, 0666)

			}

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

