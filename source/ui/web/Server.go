package web

import "exocomp/engine"
import routes_parameters "exocomp/ui/web/routes/parameters"
import routes_session "exocomp/ui/web/routes/session"
import "exocomp/types"
import utils_http "exocomp/utils/http"
import "embed"
import "fmt"
import "net/http"
import "io/fs"
import net_url "net/url"
import "os"
import "os/signal"
import "syscall"
import "time"

//go:embed public/*
var embed_fs embed.FS

type Server struct {
	Session *engine.Session
	URL     *net_url.URL
	handler *utils_http.Handler
}

func NewServer(agent *types.Agent, config *types.Config) *Server {

	session, _ := engine.NewSessionWithRecovery(agent, config)

	session.InstallTools()
	session.InstallAdapters()

	agents := session.Recovery.RestoreAgents()

	session.InstallAgents(agents)
	session.BindAgents(session.SwitchRole)

	url, _ := net_url.Parse("http://localhost:3000/")

	return &Server{
		Session: session,
		URL:     url,
		handler: utils_http.NewHandler(
			http.NotFoundHandler(),
		),
	}

}

func RestoreServer(session *engine.Session, agents []*types.Agent) *Server {

	session.InstallTools()
	session.InstallAdapters()

	session.InstallAgents(agents)
	session.BindAgents(session.SwitchRole)

	url, _ := net_url.Parse("http://localhost:3000/")

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
		server.Session.Backup()
	}

}

func (server *Server) EnableHotReload() error {

	dir_fs := os.DirFS("ui/web")
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

	if server.Session != nil && server.Session.Config != nil && server.Session.Config.GetPrompt() != "" {

		go func() {

			err := server.Session.Init()

			if err != nil {
				server.Session.Console.Error(err.Error())
			}

		}()

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
			signals <- syscall.SIGABRT
		} else {
			signals <- syscall.SIGTERM
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

	if os.Getenv("EXOCOMP_DEV") != "" {
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

	http.HandleFunc("/api/session/bugs", func(response http.ResponseWriter, request *http.Request) {
		routes_session.Bugs(server.Session, request, response)
	})

	http.HandleFunc("/api/session/changelog", func(response http.ResponseWriter, request *http.Request) {
		routes_session.Changelog(server.Session, request, response)
	})

	http.HandleFunc("/api/session/requirements", func(response http.ResponseWriter, request *http.Request) {
		routes_session.Requirements(server.Session, request, response)
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

	return http.ListenAndServe(":"+server.URL.Port(), nil)

}
