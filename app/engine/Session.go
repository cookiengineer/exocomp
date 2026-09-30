//go:build wasm

package engine

import app_schemas "exocomp-app/schemas"
import app_types "exocomp-app/types"
import "github.com/cookiengineer/gooey/components/app"
import "encoding/json"

type Session struct {
	Agent   string
	Agents  map[string]*app_types.Agent
	Client  *app.Client
	Config  *app_types.Config
	Console *app_types.Console
	Tools   []app_schemas.Tool
	Waiting bool
}

func NewSession(config *app_types.Config, client *app.Client) *Session {

	session := Session{
		Agent:   "",
		Agents:  make(map[string]*app_types.Agent),
		Config:  config,
		Console: app_types.NewConsole(),
		Tools:   make([]app_schemas.Tool, 0),
		Waiting: false,
		Client:  client,
	}

	return &session

}

// CallTool mirrors `Session.prototype.CallTool` and blocks until the request
// finished. Callers must run it inside a goroutine.
func (session *Session) CallTool(name string, method string, args map[string]any) bool {

	if session.Waiting == true {
		return false
	}

	session.Waiting = true

	arguments, err0 := json.Marshal(args)

	if err0 != nil {
		session.Waiting = false
		return false
	}

	payload, err1 := json.Marshal(app_schemas.ToolCall{
		Type: "function",
		Function: app_schemas.ToolCallFunction{
			Name:         name,
			ArgumentsRaw: arguments,
		},
	})

	if err1 != nil {
		session.Waiting = false
		return false
	}

	response, err2 := session.Client.Create("/api/session/calltool", payload)

	session.Waiting = false

	if err2 != nil || response == nil {
		return false
	}

	return response.OK

}

func (session *Session) GetAgent(name string) *app_types.Agent {

	if name == "" {
		return session.Agents[session.Agent]
	}

	agent, ok := session.Agents[name]

	if ok == true {
		return agent
	}

	return nil

}

func (session *Session) GetAgents() map[string]*app_types.Agent {
	return session.Agents
}

func (session *Session) GetMessages(from int) []*app_schemas.Message {

	result := make([]*app_schemas.Message, 0)

	agent := session.GetAgent(session.Agent)

	if agent != nil {

		if from < 0 {
			from = 0
		}

		for m := from; m < len(agent.Messages); m++ {
			result = append(result, agent.Messages[m])
		}

	}

	return result

}

func (session *Session) GetToolNames() []string {

	result := make([]string, 0)

	for _, tool := range session.Tools {
		result = append(result, tool.Function.Name)
	}

	return result

}

func (session *Session) GetToolSchema(name string) *app_schemas.Tool {

	for t := 0; t < len(session.Tools); t++ {

		if session.Tools[t].Function.Name == name {
			return &session.Tools[t]
		}

	}

	return nil

}

func (session *Session) ReceiveAgent(agent *app_types.Agent) bool {

	if agent == nil {
		return false
	}

	session.Agents[agent.Name] = agent

	if session.Agent == "" && agent.Name == session.Config.Name {
		session.Agent = agent.Name
	}

	return true

}

// SendChatRequest mirrors `Session.prototype.SendChatRequest` and blocks until
// the request finished. Callers must run it inside a goroutine.
func (session *Session) SendChatRequest(message *app_schemas.Message) bool {

	if session.Waiting == true {
		return false
	}

	session.Waiting = true

	payload, err0 := json.Marshal(message)

	if err0 != nil {
		session.Waiting = false
		return false
	}

	response, err1 := session.Client.Create("/api/session/sendchatrequest", payload)

	session.Waiting = false

	if err1 != nil || response == nil {
		return false
	}

	return response.OK

}

func (session *Session) SetAgent(agent string) bool {

	if agent == "" {
		return false
	}

	session.Agent = agent

	return true

}

func (session *Session) Update() {
	session.UpdateAgents()
}

// UpdateAgents mirrors `Session.prototype.UpdateAgents` and blocks until the
// request finished. Callers must run it inside a goroutine.
func (session *Session) UpdateAgents() {

	response, err := session.Client.Read("/api/session/agents")

	if err != nil || response == nil {
		return
	}

	agents := make([]*app_types.Agent, 0)
	err = json.Unmarshal(response.Body, &agents)

	if err != nil {
		return
	}

	for _, agent := range agents {
		session.ReceiveAgent(agent)
	}

}

// UpdateTools mirrors `Session.prototype.UpdateTools` and blocks until the
// request finished. Callers must run it inside a goroutine.
func (session *Session) UpdateTools() {

	response, err := session.Client.Read("/api/session/tools")

	if err != nil || response == nil {
		return
	}

	tools := make([]app_schemas.Tool, 0)
	err = json.Unmarshal(response.Body, &tools)

	if err != nil {
		return
	}

	session.Tools = session.Tools[:0]

	for _, tool := range tools {
		session.Tools = append(session.Tools, tool)
	}

}
