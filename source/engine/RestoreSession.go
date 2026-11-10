package engine

import "exocomp/types"
import utils_chat "exocomp/utils/chat"
import net_http "net/http"
import "os"
import "path/filepath"
import "sync"

func RestoreSession(playground string, backup Session) *Session {

	client := &net_http.Client{
		Timeout: utils_chat.RetryRequestTimeout,
	}

	session := &Session{
		Agent:    backup.Agent,
		Config:   backup.Config,
		Console:  types.NewConsole(os.Stdout, os.Stderr, 0),
		Recovery: NewRecovery(playground),
		Waiting:  false,
		client:   client,
		mutex:    &sync.RWMutex{},
		adapters: make(map[string]types.Adapter),
		tools:    make(map[string]types.Tool),
		changes:  make(chan struct{}, 1),
	}

	if backup.Console != nil {
		session.Console.Messages = backup.Console.Messages
	}

	if backup.Config != nil {

		if backup.Config.Playground == backup.Config.Sandbox {

			session.Config.Playground = playground
			session.Config.Sandbox = playground

		} else {

			relative_sandbox, err3 := filepath.Rel(backup.Config.Playground, backup.Config.Sandbox)

			if err3 == nil {

				session.Config.Playground = playground
				session.Config.Sandbox = filepath.Join(playground, relative_sandbox)

			}

		}

	}

	if session.Agent != nil && session.Config != nil {
		session.Agent.Sandbox = session.Config.Sandbox
	}

	return session

}
