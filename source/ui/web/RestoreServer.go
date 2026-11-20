package web

import "exocomp/engine"
import "exocomp/types"
import utils_http "exocomp/utils/http"
import "net/http"
import net_url "net/url"

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
