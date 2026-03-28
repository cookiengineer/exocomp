package ollama

import "exocomp/agents"
import "exocomp/config"
import "exocomp/tools"
import "fmt"
import "net/http"
import "strings"
import "sync"

type Session struct {
	Agent    *agents.Agent
	Config   *config.Config
	Client   *http.Client
	Messages []*Message
	Tools    []*tools.Tool
	Waiting  bool
	mutex    *sync.Mutex
}

func NewSession(agent *agents.Agent, config *config.Config) (*Session, error) {

	session := &Session{
		Agent:    agent,
		Config:   config,
		Client:   &http.Client{},
		Messages: make([]*Message, 0),
		Waiting:  false,
		mutex:    &sync.Mutex{},
	}

	system_prompts := make([]string, 0)
	system_tools   := tools.ToSchema(config.Tools)

	if agent != nil {
		system_prompts = append(system_prompts, agent.GetPrompt())
	}

	if config != nil {
		system_prompts = append(system_prompts, config.GetPrompt(system_tools))
	}

	session.mutex.Lock()
	session.Messages = append(session.Messages, &Message{
		Role:    "system",
		Content: strings.Join(system_prompts, "\n"),
		Tools:   system_tools,
	})
	session.mutex.Unlock()

	err := sendChatRequest(session)

	return session, err

}

func (session *Session) LastMessage() *Message {

	if len(session.Messages) > 0 {
		return session.Messages[len(session.Messages) - 1]
	} else {
		return nil
	}

}

func (session *Session) Query(message Message) error {

	if session.Waiting == false {

		session.mutex.Lock()
		session.Messages = append(session.Messages, &message)
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

