package web

import "exocomp/agents"
import "exocomp/types"
import "exocomp/ui/web/routes"
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

	// TODO: Remove this when finished
	embed_fs := os.DirFS("ui/web")

	fsys, _ := fs.Sub(embed_fs, "public")
	fsrv    := http.FileServer(http.FS(fsys))

	http.Handle("/", fsrv)
	// func(response http.ResponseWriter, request *http.Request) {
	// 	fsrv.ServeHTTP(response, request)
	// })

	http.HandleFunc("/api/parameters/agents", func(response http.ResponseWriter, request *http.Request) {
		routes.AgentParameters(server.Session, request, response)
	})

	http.HandleFunc("/api/parameters/models", func(response http.ResponseWriter, request *http.Request) {
		routes.ModelParameters(server.Session, request, response)
	})

	http.HandleFunc("/api/console", func(response http.ResponseWriter, request *http.Request) {
		routes.Console(server.Session, request, response)
	})

	http.HandleFunc("/api/agents", func(response http.ResponseWriter, request *http.Request) {
		routes.Agents(server.Session, request, response)
	})

	http.HandleFunc("/api/messages", func(response http.ResponseWriter, request *http.Request) {
		routes.Messages(server.Session, request, response)
	})

	http.HandleFunc("/api/messages/send", func(response http.ResponseWriter, request *http.Request) {
		routes.SendMessage(server.Session, request, response)
	})

	http.HandleFunc("/api/settings/agent", func(response http.ResponseWriter, request *http.Request) {
		routes.AgentSettings(server.Session, request, response)
	})

	err := http.ListenAndServe(":" + server.URL.Port(), nil)

	if err == nil {
		return true
	}

	return false

}
