package session

import "exocomp/ui/web/handlers"
import "exocomp/tools"
import "exocomp/types"
import "encoding/json"
import "io"
import "net/http"
import "strconv"

func Answer(session *types.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodPost {

		content_type := request.Header.Get("Content-Type")

		if content_type == "application/json" {

			request_payload, err0 := io.ReadAll(request.Body)

			if err0 == nil {

				answer := struct {
					ID     string `json:"id"`
					Answer string `json:"answer"`
				}{}

				err1 := json.Unmarshal(request_payload, &answer)

				if err1 == nil {

					tool := session.GetTool("questions.Ask")

					if tool != nil {

						questions_tool, ok := tool.(*tools.Questions)

						if ok == true {

							err2 := questions_tool.Answer(answer.ID, answer.Answer)

							if err2 == nil {

								content, _ := question_tool.GetContent(answer.ID)
								response_payload, _ := json.MarshalIndent(content, "", "\t")

								response.Header().Set("Content-Type", "application/json")
								response.Header().Set("Content-Length", strconv.Itoa(len(response_payload)))
								response.WriteHeader(http.StatusOK)
								response.Write(response_payload)

							} else {
								handlers.BadRequest(session, err2, request, response)
							}

						} else {
							handlers.NotFound(session, request, response)
						}

					} else {
						handlers.NotFound(session, request, response)
					}

				} else {
					handlers.BadRequest(session, err1, request, response)
				}

			} else {
				handlers.UnsupportedMediaType(session, request, response)
			}

		} else {
			handlers.UnsupportedMediaType(session, request, response)
		}

	} else {
		handlers.MethodNotAllowed(session, request, response)
	}

}
