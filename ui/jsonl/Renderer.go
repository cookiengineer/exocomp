package jsonl

import "exocomp/schemas"
import "exocomp/types"
import "encoding/json"
import "fmt"
import "os"
import "sync"

type Renderer struct {
	Session   *types.Session
	mutex     *sync.RWMutex
	rendered  int
}

func NewRenderer(session *types.Session) *Renderer {

	return &Renderer{
		Session:  session,
		mutex:    &sync.RWMutex{},
		rendered: 0,
	}

}

func (renderer *Renderer) Destroy() {

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
		}

		bytes, err := json.Marshal(dummy)

		if err == nil {

			fmt.Fprintf(os.Stdout, "%s\n", string(bytes))
			os.Stdout.Sync()

		}

	}

}

