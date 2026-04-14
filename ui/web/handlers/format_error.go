package handlers

import "fmt"
import "net/http"
import "strings"

func format_error(request *http.Request, message string) (string, []byte) {

	content_type := "text/html; charset=utf-8"
	payload := []byte(fmt.Sprintf(`
<!DOCTYPE html>
<html>
	<head>
		<style>
		body {
			margin: 0px;
			padding: 0px;
			color: #cc0000;
			background: #000000;
			font-size: 32px;
			line-height: 100vh;
			text-align: center;
			vertical-align: middle;
			overflow: hidden;
			animation: thedarkside 2000ms ease-in-out infinite;
		}
		@keyframes thedarkside {
			0%%, 100%% {
				text-shadow:
				  0 0 2px #660000,
				  0 0 6px #990000,
				  0 0 10px #cc0000;
			}
			50%% {
				text-shadow:
				  0 0 4px #990000,
				  0 0 10px #cc0000,
				  0 0 20px rgba(204, 0, 0, 0.7);
			}
		}
		</style>
	</head>
	<body><div>%s</div></body>
</html>
`, message))

	if strings.Contains(request.URL.Path, "/api/") {
		content_type = "application/json; charset=utf-8"
		payload = []byte(fmt.Sprintf("{\"error\": \"%s\"}", message))
	} else if strings.HasSuffix(request.URL.Path, ".json") {
		content_type = "application/json; charset=utf-8"
		payload = []byte(fmt.Sprintf("{\"error\": \"%s\"}", message))
	}

	return content_type, payload

}
