package web

import routes_parameters "exocomp/ui/web/routes/parameters"
import routes_session "exocomp/ui/web/routes/session"
import "exocomp/tools"
import "exocomp/types"
import utils_http "exocomp/utils/http"
import "embed"
import "fmt"
import "net/http"
import "io/fs"
import net_url "net/url"
import "os"
import "os/signal"
import "strings"
import "syscall"
import "time"

//go:embed public/*
var embed_fs embed.FS

type Server struct {
	Session *types.Session
	URL     *net_url.URL
	handler *utils_http.Handler
}

func NewServer(agent *types.Agent, config *types.Config) *Server {

	var session *types.Session = nil

	recovery := types.NewRecovery(config.Playground)

	if recovery.HasBackup() {

		session = recovery.RestoreSession()

		if session == nil {
			session = types.NewSession(agent, config)
		}

	} else {
		session = types.NewSession(agent, config)
	}

	url, _ := net_url.Parse("http://localhost:3000/")

	if len(agent.AllowedTools) > 0 {

		tool_schemas, tools := tools.Toolset(
			config.Playground,
			config.Sandbox,
			config.Model,
			config.URL,
			config.Debug,
			agent.AllowedPrograms,
			agent.AllowedTools,
		)

		for name, tool := range tools {
			session.SetTool(name, tool, tool_schemas[name])
		}

	}

	tool := session.GetTool("agents.List")

	if tool != nil {

		agent_tool, ok := tool.(*tools.Agents)

		if ok == true {

			agents := recovery.RestoreAgents()

			if len(agents) > 0 {

				for _, agent := range agents {
					agent_tool.SetAgent(agent)
				}

			}

		}

	}

	return &Server{
		Session: session,
		URL:     url,
		handler: utils_http.NewHandler(
			http.NotFoundHandler(),
		),
	}

}

func (server *Server) Destroy() {

	if server.Session != nil {

		server.Session.Recovery.BackupSession(server.Session)

		tool := server.Session.GetTool("agents.List")

		if tool != nil {

			agent_tool, ok := tool.(*tools.Agents)

			if ok == true {

				agent_names := agent_tool.GetNames()

				for _, name := range agent_names {

					agent := agent_tool.GetAgent(name)

					if agent != nil {
						server.Session.Recovery.BackupAgent(agent)
					}

				}

			}

		}

	}

}

func (server *Server) EnableHotReload() error {

	dir_fs     := os.DirFS("ui/web")
	fsys, err0 := fs.Sub(dir_fs, "public")

	if err0 == nil {

		server.handler.Set(http.FileServer(http.FS(fsys)))

		return nil

	} else {
		return err0
	}

}

func (server *Server) Init() {

	fsys, err0 := fs.Sub(embed_fs, "public")

	if err0 == nil {
		server.handler.Set(http.FileServer(http.FS(fsys)))
	} else {
		panic(err0)
	}

	signals := make(chan os.Signal, 1)

	signal.Notify(
		signals,
		syscall.SIGABRT,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go func() {

		err := server.Listen()

		if err != nil {
			signals<-syscall.SIGABRT
		} else {
			signals<-syscall.SIGTERM
		}

	}()

	select {
	case sig := <-signals:

		switch sig {
		case syscall.SIGABRT:

			server.Destroy()
			fmt.Fprintf(os.Stdout, "Received signal: %s\n", "SIGABRT")

			time.Sleep(1 * time.Second)
			os.Exit(1)

		case syscall.SIGINT:

			server.Destroy()
			fmt.Fprintf(os.Stdout, "Received signal: %s\n", "SIGINT")

			time.Sleep(1 * time.Second)
			os.Exit(0)

		case syscall.SIGTERM:

			server.Destroy()
			fmt.Fprintf(os.Stdout, "Received signal: %s\n", "SIGTERM")

			time.Sleep(1 * time.Second)
			os.Exit(0)

		default:

			server.Destroy()
			fmt.Printf("Received signal: %s\n", sig.String())

			time.Sleep(1 * time.Second)
			os.Exit(0)

		}

	}

}

func (server *Server) Listen() error {

	cwd, _ := os.Getwd()

	if strings.HasSuffix(cwd, "exocomp/source") {
		server.EnableHotReload()
	}

	// NOTE: handler is an atomic value
	http.Handle("/", server.handler)

	// CLI Parameters
	http.HandleFunc("/api/parameters/roles", func(response http.ResponseWriter, request *http.Request) {
		routes_parameters.Roles(server.Session, request, response)
	})

	http.HandleFunc("/api/parameters/models", func(response http.ResponseWriter, request *http.Request) {
		routes_parameters.Models(server.Session, request, response)
	})

	// Session
	http.HandleFunc("/api/session/config", func(response http.ResponseWriter, request *http.Request) {
		routes_session.Config(server.Session, request, response)
	})

	http.HandleFunc("/api/session/config/{name}", func(response http.ResponseWriter, request *http.Request) {
		routes_session.AgentConfig(server.Session, request, response)
	})

	http.HandleFunc("/api/session/agent", func(response http.ResponseWriter, request *http.Request) {
		routes_session.Agent(server.Session, request, response)
	})

	http.HandleFunc("/api/session/agents", func(response http.ResponseWriter, request *http.Request) {
		routes_session.Agents(server.Session, request, response)
	})

	http.HandleFunc("/api/session/console", func(response http.ResponseWriter, request *http.Request) {
		routes_session.Console(server.Session, request, response)
	})

	http.HandleFunc("/api/session/tools", func(response http.ResponseWriter, request *http.Request) {
		routes_session.Tools(server.Session, request, response)
	})


	// Session Interaction
	http.HandleFunc("/api/session/calltool", func(response http.ResponseWriter, request *http.Request) {
		routes_session.CallTool(server.Session, request, response)
	})

	http.HandleFunc("/api/session/sendchatrequest", func(response http.ResponseWriter, request *http.Request) {
		routes_session.SendChatRequest(server.Session, request, response)
	})


	return http.ListenAndServe(":" + server.URL.Port(), nil)

}

