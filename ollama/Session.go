package ollama

import "bytes"
import "encoding/json"
import "fmt"
import "io"
import "net/http"
import "exocomp/config"

type Session struct {
	agent   *agents.Agent
	config  *config.Config
	client  *http.Client
	history []*Message
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string     `json:"model"`
	Messages []*Message `json:"messages"`
	Stream   bool       `json:"stream"`
}

type ChatResponse struct {
	Message Message `json:"message"`
}

func NewSession(agent *agents.Agent, config *config.Config) (*Session, error) {

	session := &Session{
		agent:   agent,
		config:  config,
		client:  &http.Client{},
		history: make([]*Message, 0),
	}

	system_prompts := make([]string, 0)

	if agent != nil {
		system_prompts = append(system_prompts, agent.GetPrompt())
	}

	if config != nil {
		system_prompts = append(system_prompts, config.GetPrompt())
	}

	session.history = append(session.history, &Message{
		Role:    "system",
		Content: strings.Join(system_prompts, "\n"),
	})

	_, err := session.send()

	if err == nil {
		return session, nil
	} else {
		return nil, err
	}

}

func (session *Session) Query(message string) (string, error) {

	session.history = append(session.history, &Message{
		Role:    "user",
		Content: message,
	})

	response, err1 := session.send()

	if response != nil && err1 == nil {

		session.history = append(session.history, response)

		gadget := config.ParseGadget(response.Content)

		if gadget != nil && session.config.IsAllowedGadget(gadget.Type.String()) {

			result, err2 := gadget.Execute(session.config)

			if err2 != nil {
				result = fmt.Sprintf("Gadget Error: %s", err2.Error())
			}

			session.history = append(session.history, &Message{
				Role:    "user",
				Content: formatGadgetResult(result),
			})

			return response.Content, nil

		} else {
			return response.Content, nil
		}

	} else {
		return "", err1
	}

}

func (session *Session) send() (*Message, error) {

	request_payload, err0 := json.Marshal(ChatRequest{
		Model:    session.config.Model,
		Messages: session.history,
		Stream:   false,
	})

	if err0 == nil {

		endpoint := session.config.ResolvePath("/api/chat")

		if session.config.Verbose == true {
			fmt.Println("POST", endpoint.String())
		}

		response, err1 := session.client.Post(
			endpoint.String(),
			"application/json",
			bytes.NewReader(request_payload),
		)

		if err1 == nil {

			response_payload, err2 := io.ReadAll(response.Body)

			if err2 == nil {

				var response ChatResponse

				err3 := json.Unmarshal(response_payload, &response)

				if err3 == nil {
					return &response.Message, nil
				} else {
					return nil, err3
				}

			} else {
				return nil, err2
			}

		} else {
			return nil, err1
		}

	} else {
		return nil, err0
	}

}
