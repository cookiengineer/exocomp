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

	http.HandleFunc("/api/init", func(response http.ResponseWriter, request *http.Request) {

		// TODO: Call server.Session.Init()

	})

	http.HandleFunc("/api/chat", func(response http.ResponseWriter, request *http.Request) {
		routes.Chat(server.Session, response, request)
	})

	http.HandleFunc("/api/models", func(response http.ResponseWriter, request *http.Request) {
		routes.Models(server.Session, response, request)
	})

	http.HandleFunc("/api/models", func(response http.ResponseWriter, request *http.Request) {

		// TODO: session.QueryModels()

	})

	err := http.ListenAndServe(":" + server.URL.Port(), nil)

	if err == nil {
		return true
	}

	return false

}
