package session

import "exocomp/ui/web/handlers"
import "exocomp/schemas"
import "exocomp/types"
import "encoding/json"
import "net/http"
import "io"
import "strconv"

func CallTool(session *types.Session, request *http.Request, response http.ResponseWriter) {

	if request.Method == http.MethodPost {

		content_type := request.Header.Get("Content-Type")

		if content_type == "application/json" {

			request_payload, err0 := io.ReadAll(request.Body)

			if err0 == nil {

				tool_call := schemas.ToolCall{}
				err1      := json.Unmarshal(request_payload, &tool_call)

				if err1 == nil {

					name,      err2 := tool_call.Function.ToName()
					method,    err3 := tool_call.Function.ToMethod()
					arguments, err4 := tool_call.Function.ToArguments()

					if err2 == nil {

						if err3 == nil {

							if err4 == nil {

								err5 := session.CallTool(name, method, arguments)

								if err5 == nil {

									response_payload, _ := json.MarshalIndent(map[string]string{
										"status": "OK",
									}, "", "\t")

									response.Header().Set("Content-Type", "application/json")
									response.Header().Set("Content-Length", strconv.Itoa(len(response_payload)))
									response.WriteHeader(http.StatusOK)
									response.Write(response_payload)

								} else {
									handlers.Unauthorized(session, err5, request, response)
								}

							} else {
								handlers.BadRequest(session, err4, request, response)
							}

						} else {
							handlers.BadRequest(session, err3, request, response)
						}

					} else {
						handlers.BadRequest(session, err2, request, response)
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
