package ollama

import "exocomp/agents"
import "exocomp/config"
import "exocomp/schemas"
import "exocomp/tools"
import "fmt"
import "net/http"
import "strings"
import "sync"

type Session struct {
	Agent    *agents.Agent
	Config   *config.Config
	Client   *http.Client
	Messages []schemas.Message
	Tools    []schemas.Tool
	Waiting  bool
	mutex    *sync.Mutex
}

func NewSession(agent *agents.Agent, config *config.Config) (*Session, error) {

	session := &Session{
		Agent:    agent,
		Config:   config,
		Client:   &http.Client{},
		Messages: make([]schemas.Message, 0),
		Tools:    make([]schemas.Tool, 0),
		Waiting:  false,
		mutex:    &sync.Mutex{},
	}

	if len(config.Tools) > 0 {
		session.Tools = tools.EncodeSchema(config.Tools)
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

	err := sendChatRequest(session)

	return session, err

}

func (session *Session) Query(message schemas.Message) error {

	if session.Waiting == false {

		session.mutex.Lock()
		session.Messages = append(session.Messages, message)
		session.Waiting = true
		session.mutex.Unlock()

		err := sendChatRequest(session)

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

		if name == "bugs" {

			// TODO
			return nil

		} else if name == "features" {

			// TODO
			return nil

		} else if name == "files" {

			return tools.Tool(tools.NewFiles(session.Config.Agent, session.Config.Sandbox))

		} else if name == "notes" {

			return tools.Tool(tools.NewNotes(session.Config.Agent, session.Config.Sandbox))

		} else if name == "programs" {

			return tools.Tool(tools.NewPrograms(session.Config.Agent, session.Config.Sandbox, session.Config.Programs))

		} else {
			return nil
		}

	} else {
		return nil
	}

}
