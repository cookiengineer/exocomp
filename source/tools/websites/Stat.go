package websites

import "fmt"
import "io"
import net_http "net/http"
import "sort"
import "strings"

func (tool *Websites) Stat(url_str string, user_agent string) (string, error) {

	url, err0 := parseWebsiteURL(url_str)

	if err0 == nil {

		useragent, err1 := getWebsiteUserAgent(user_agent)

		if err1 == nil {

			if tool.canConnect(url.Hostname()) == false {
				return "", fmt.Errorf("websites.Stat: Missing network.Connect permission for \"%s\". Request it with humans.AskPermission.", url.Hostname())
			}

			var response *net_http.Response = nil

			request_head, err2 := net_http.NewRequest(net_http.MethodHead, url.String(), nil)

			if err2 == nil {

				request_head.Header = useragent.Header()

				response_head, err3 := tool.client.Do(request_head)

				if err3 == nil {

					io.Copy(io.Discard, response_head.Body)
					response_head.Body.Close()

					// NOTE: Falls through for response handling
					response = response_head

				} else {

					// Some servers reject HEAD, fall back to a body-less GET
					request_get, err4 := net_http.NewRequest(net_http.MethodGet, url.String(), nil)

					if err4 == nil {

						request_get.Header = useragent.Header()

						response_get, err5 := tool.client.Do(request_get)

						if err5 == nil {

							io.Copy(io.Discard, response_get.Body)
							response_get.Body.Close()

							// NOTE: Falls through for response handling
							response = response_get

						} else {
							return "", fmt.Errorf("websites.Stat: URL \"%s\" is unreachable.", url.String())
						}

					} else {
						return "", fmt.Errorf("websites.Stat: Invalid URL \"%s\".", url.String())
					}

				}

			} else {
				return "", fmt.Errorf("websites.Stat: Invalid URL \"%s\".", url.String())
			}

			if response != nil {

				lines := make([]string, 0)
				lines = append(lines, fmt.Sprintf("websites.Stat: %s (%s)", url.String(), response.Status))

				keys := make([]string, 0)

				for key, _ := range response.Header {
					keys = append(keys, key)
				}

				sort.Strings(keys)

				for _, key := range keys {
					lines = append(lines, fmt.Sprintf("%s: %s", key, strings.Join(response.Header[key], ", ")))
				}

				return strings.Join(lines, "\n"), nil

			} else {
				return "", fmt.Errorf("websites.Stat: URL \"%s\" gave no HTTP Response. Maybe a connection error?", url.String())
			}

		} else {
			return "", fmt.Errorf("websites.Stat: %s", err1.Error())
		}

	} else {
		return "", fmt.Errorf("websites.Stat: %s", err0.Error())
	}

}
