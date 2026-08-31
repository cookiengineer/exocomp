package jsonl

import "exocomp/engine"
import "exocomp/schemas"
import "encoding/json"
import "fmt"
import "os"
import "sync"
import "time"

type Renderer struct {
	Session   *engine.Session
	mutex     *sync.RWMutex
	rendered  int
}

func NewRenderer(session *engine.Session) *Renderer {

	return &Renderer{
		Session:  session,
		mutex:    &sync.RWMutex{},
		rendered: 0,
	}

}

func (renderer *Renderer) Destroy() {

}

func (renderer *Renderer) Flushed() bool {

	renderer.mutex.RLock()
	defer renderer.mutex.RUnlock()

	if renderer.Session != nil && renderer.Session.Agent != nil {
		return renderer.rendered >= len(renderer.Session.Agent.Messages)
	}

	return true

}

func (renderer *Renderer) RenderLoop() {

	os.Stdout.Sync()

	for {

		renderer.mutex.RLock()
		from := renderer.rendered
		renderer.mutex.RUnlock()

		messages := renderer.Session.GetMessages(from)

		if len(messages) > 0 {

			renderer.RenderMessages(messages)

			renderer.mutex.Lock()
			renderer.rendered += len(messages)
			renderer.mutex.Unlock()

		} else {

			time.Sleep(100 * time.Millisecond)
			continue

		}

	}

}

func (renderer *Renderer) RenderMessages(messages []*schemas.Message) {

	for _, message := range messages {

		if message == nil {
			continue
		}

		dummy := schemas.Message{
			Role:    message.Role,
			Content: message.Content,
			Created: message.Created,
		}

		bytes, err := json.Marshal(dummy)

		if err == nil {

			fmt.Fprintf(os.Stdout, "schemas.Message:%s\n", string(bytes))
			os.Stdout.Sync()

		}

	}

}

