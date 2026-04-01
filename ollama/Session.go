package ollama

import "exocomp/agents"
import "exocomp/schemas"
import "exocomp/tools"
import "exocomp/types"
import "fmt"
import "net/http"
import "strings"
import "sync"

type Session struct {
	Agent    *agents.Agent
	Config   *types.Config
	Client   *http.Client
	Messages []schemas.Message
	Tools    []schemas.Tool
	Waiting  bool
	mutex    *sync.Mutex
}

func NewSession(agent *agents.Agent, config *types.Config) *Session {

	session := &Session{
		Agent:    agent,
		Config:   config,
		Client:   &http.Client{},
		Messages: make([]schemas.Message, 0),
		Tools:    make([]schemas.Tool, 0),
		Waiting:  false,
		mutex:    &sync.Mutex{},
	}

	if len(agent.Tools) > 0 {
		session.Tools = tools.EncodeSchema(agent.Tools)
	}

	system_prompts := make([]string, 0)

	if agent != nil {
		system_prompts = append(system_prompts, agent.GetPrompt())
	}

	if config != nil {
		system_prompts = append(system_prompts, config.GetPrompt())
	}

	session.mutex.Lock()
	session.Messages = append(session.Messages, schemas.Message{
		Role:    "system",
		Content: strings.Join(system_prompts, "\n"),
	})
	session.mutex.Unlock()

	return session

}

func (session *Session) Init() error {

	if len(session.Messages) > 0 {
		return SendChatRequest(session)
	}

	return fmt.Errorf("Session is empty, waiting for LLM system prompt ...")

}

func (session *Session) Query(message schemas.Message) error {

	if session.Waiting == false {

		session.mutex.Lock()
		session.Messages = append(session.Messages, message)
		session.Waiting = true
		session.mutex.Unlock()

		err := SendChatRequest(session)

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

func (session *Session) GetTool(name string) tools.Tool {

	allowed := false

	for _, tool := range session.Tools {

		if strings.HasPrefix(tool.Function.Name, name + ".") {
			allowed = true
			break
		}

	}

	if allowed == true {

		if name == "agents" {

			return tools.Tool(tools.NewAgents(session.Config.Agent, session.Config.Sandbox))

		} else if name == "bugs" {

			return tools.Tool(tools.NewBugs(session.Config.Agent, session.Config.Sandbox))

		} else if name == "changelog" {

			return tools.Tool(tools.NewChangelog(session.Config.Agent, session.Config.Sandbox))

		} else if name == "files" {

			return tools.Tool(tools.NewFiles(session.Config.Agent, session.Config.Sandbox))

		} else if name == "programs" {

			return tools.Tool(tools.NewPrograms(session.Config.Agent, session.Config.Sandbox, session.Agent.Programs))

		} else if name == "requirements" {

			return tools.Tool(tools.NewRequirements(session.Config.Agent, session.Config.Sandbox))

		} else {
			return nil
		}

	} else {
		return nil
	}

}
