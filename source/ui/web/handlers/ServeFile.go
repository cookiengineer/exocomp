package handlers

import "exocomp/types"
import "bytes"
import "compress/gzip"
import "net/http"
import "mime"
import "path"
import "strconv"
import "strings"

func ServeFile(session *types.Session, request *http.Request, response http.ResponseWriter, data []byte) {

	if request.Method == http.MethodGet {

		ext := path.Ext(request.URL.Path)
		content_type := mime.TypeByExtension(ext)

		if content_type == "" {
			content_type = "application/octet-stream"
		}

		if ext == ".html" {

			policies := []string{

				// Assets
				"default-src 'self'",
				"script-src 'self'",
				"script-src-elem 'self'",
				"font-src 'self'",
				"img-src 'self'",
				"style-src 'self'",

				// Disallow Trackers
				"prefetch-src 'none'",
				"worker-src 'self'",

				// Allow XHR/Fetch to other websites
				"media-src 'self'",
				"frame-src * 'self'",
				"connect-src * 'self'",
			}

			for _, policy := range policies {
				response.Header().Set("Content-Security-Policy", policy)
			}

		}

		accept_encoding := request.Header.Get("Accept-Encoding")

		if strings.Contains(accept_encoding, "gzip") {

			var buffer bytes.Buffer

			writer := gzip.NewWriter(&buffer)

			_, err1 := writer.Write(data);

			if err1 == nil {

				err2 := writer.Close()

				if err2 == nil {

					encoded := buffer.Bytes()

					response.Header().Set("Content-Type", content_type)
					response.Header().Set("Content-Encoding", "gzip")
					response.Header().Set("Content-Length", strconv.Itoa(len(encoded)))
					response.Write(encoded)

				} else {
					InternalServerError(session, err2, request, response)
				}

			} else {
				InternalServerError(session, err1, request, response)
			}

		} else {

			response.Header().Set("Content-Type", content_type)
			response.Header().Set("Content-Encoding", "identity")
			response.Header().Set("Content-Length", strconv.Itoa(len(data)))
			response.Write(data)

		}

	} else {
		MethodNotAllowed(session, request, response)
	}

}
