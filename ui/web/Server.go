package web

import "exocomp/agents"
import "exocomp/types"
import "exocomp/ui/web/routes"
import routes_agents "exocomp/ui/web/routes/agents"
import routes_parameters "exocomp/ui/web/routes/parameters"
import routes_session "exocomp/ui/web/routes/session"
import "embed"
import "net/http"
import "io/fs"
import net_url "net/url"
import "os"

//go:embed public/*
var embed_fs embed.FS

type Server struct {
	Session  *types.Session
	URL      *net_url.URL
}

func NewServer(agent *agents.Agent, config *types.Config) *Server {

	session := types.NewSession(agent, config)
	url, _ := net_url.Parse("http://localhost:3000/")

	return &Server{
		Session: session,
		URL:     url,
	}

}

func (server *Server) Init() bool {

	if server.Session.Config.Debug == true {

		dir_fs    := os.DirFS("ui/web")
		fsys, err := fs.Sub(dir_fs, "public")

		if err == nil {
			fsrv := http.FileServer(http.FS(fsys))
			http.Handle("/", fsrv)
		}

	} else {

		fsys, err := fs.Sub(embed_fs, "public")

		if err == nil {
			fsrv := http.FileServer(http.FS(fsys))
			http.Handle("/", fsrv)
		}

	}

	// CLI Parameters
	http.HandleFunc("/api/parameters/agents", func(response http.ResponseWriter, request *http.Request) {
		routes_parameters.Agents(server.Session, request, response)
	})

	http.HandleFunc("/api/parameters/models", func(response http.ResponseWriter, request *http.Request) {
		routes_parameters.Models(server.Session, request, response)
	})

	// Session
	http.HandleFunc("/api/session/config", func(response http.ResponseWriter, request *http.Request) {
		routes_session.Config(server.Session, request, response)
	})

	http.HandleFunc("/api/session/console", func(response http.ResponseWriter, request *http.Request) {
		routes_session.Console(server.Session, request, response)
	})

	http.HandleFunc("/api/session/context", func(response http.ResponseWriter, request *http.Request) {
		routes_session.Context(server.Session, request, response)
	})

	http.HandleFunc("/api/session/messages", func(response http.ResponseWriter, request *http.Request) {
		routes_session.Messages(server.Session, request, response)
	})

	http.HandleFunc("/api/session/sendchatrequest", func(response http.ResponseWriter, request *http.Request) {
		routes_session.SendChatRequest(server.Session, request, response)
	})

	// Agents
	http.HandleFunc("/api/agents", func(response http.ResponseWriter, request *http.Request) {
		routes_agents.Index(server.Session, request, response)
	})

	http.HandleFunc("/api/agents/{agent}", func(response http.ResponseWriter, request *http.Request) {
		routes_agents.Agent(server.Session, request, response)
	})



	// TODO
	http.HandleFunc("/api/settings/agent", func(response http.ResponseWriter, request *http.Request) {
		routes.AgentSettings(server.Session, request, response)
	})

	err := http.ListenAndServe(":" + server.URL.Port(), nil)

	if err == nil {
		return true
	}

	return false

}
