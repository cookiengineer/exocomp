package actions

import "exocomp/engine"
import "exocomp/types"
import "exocomp/ui/tty"
import "exocomp/ui/web"
import "fmt"
import "os"
import "strings"

func Debug(session *engine.Session, agents []*types.Agent, role string) {

	client := tty.RestoreClient(session, agents)
	server := web.RestoreServer(session, agents)

	if role != "" {
		client.SetRole(role)
	}

	fmt.Fprintf(os.Stdout, "[config]:\n")
	fmt.Fprintf(os.Stdout, "| Agent:   %s | %s | %s | %.2f\n", session.Agent.Name, session.Agent.Role, session.Agent.Model, session.Agent.Temperature)
	fmt.Fprintf(os.Stdout, "| Tools:   %s\n", strings.Join(session.Agent.AllowedTools, ", "))
	fmt.Fprintf(os.Stdout, "| Sandbox: %s\n", session.Config.Sandbox)
	fmt.Fprintf(os.Stdout, "| URL:     %s\n", session.Config.URL.String())
	fmt.Fprintf(os.Stdout, "| Web:     %s\n", server.URL.String())
	fmt.Fprintf(os.Stdout, "| Debug:   %t\n", session.Config.Debug)
	fmt.Fprintf(os.Stdout, "|\n")

	for _, agent := range agents {
		fmt.Fprintf(os.Stdout, "|-> Restored Agent \"%s\" with %d messages\n", agent.Name, len(agent.Messages))
	}

	os.Stdout.Sync()


	shutdown := make(chan bool, 1)

	go func() {
		client.Init()
		shutdown<-true
	}()

	go func() {
		server.Init()
		shutdown<-true
	}()

	select {
	case <-shutdown:
		os.Exit(0)
	}

}
