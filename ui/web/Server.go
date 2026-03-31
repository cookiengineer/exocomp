package web

import "exocomp/agents"
import "exocomp/ollama"
import "exocomp/types"
import "embed"
import "net/http"
import "io/fs"
import net_url "net/url"
import "os"

//go:embed public/*
var embed_fs embed.FS

type Server struct {
	Session  *ollama.Session
	URL      *net_url.URL
}

func NewServer(agent *agents.Agent, config *types.Config) *Server {

	session, err0 := ollama.NewSession(agent, config)

	if err0 == nil {

		url, _ := net_url.Parse("http://localhost:1234/")

		return &Server{
			Session: session,
			URL:     url,
		}

	} else {
		return nil
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

	http.HandleFunc("/api/chat", func(response http.ResponseWriter, request *http.Request) {

		// TODO: Deserialize message payload
		// TODO: session.Query(message)
		// TODO: Respond with response messages

	})

	err := http.ListenAndServe(":" + server.URL.Port(), nil)

	if err == nil {
		return true
	}

	return false

}
