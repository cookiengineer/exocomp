package websites

import "fmt"
import "io"
import net_http "net/http"
import "strings"

func (tool *Websites) Fetch(url_str string, user_agent string, format string) (string, error) {

	format = strings.ToLower(strings.TrimSpace(format))

	switch format {

	case "markdown", "text", "html":
	case "":
		format = "markdown"
	default:
		format = ""

	}

	if format != "" {

		url, err0 := parseWebsiteURL(url_str)

		if err0 == nil {

			useragent, err1 := getWebsiteUserAgent(user_agent)

			if err1 == nil {

				if tool.canConnect(url.Hostname()) == false {
					return "", fmt.Errorf("websites.Fetch: Missing network.Connect permission for \"%s\". Request it with humans.AskPermission.", url.Hostname())
				}

				request, err2 := net_http.NewRequest(net_http.MethodGet, url.String(), nil)

				if err2 == nil {

					request.Header = useragent.Header()

					accept := request.Header.Get("Accept")
					accept_language := request.Header.Get("Accept-Language")

					if accept == "" {
						accept = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
					}

					if accept_language == "" {
						accept_language = "en-US,en;q=0.9"
					}

					request.Header.Set("Accept", accept)
					request.Header.Set("Accept-Language", accept_language)

					response, err3 := tool.client.Do(request)

					if err3 == nil {

						defer response.Body.Close()

						bytes, err4 := io.ReadAll(io.LimitReader(response.Body, 16*1024*1024))

						if err4 == nil {

							content := ""

							switch format {

							case "html":
								content = string(bytes)
							case "text":
								content = getWebsiteText(url, bytes)
							default:
								content = getWebsiteMarkdown(url, bytes)

							}

							return strings.Join([]string{
								fmt.Sprintf("websites.Fetch: URL \"%s\" (%s) contents", url.String(), response.Status),
								strings.TrimSpace(content),
							}, "\n"), nil

						} else {
							return "", fmt.Errorf("websites.Fetch: Cannot read HTTP response body.")
						}

					} else {
						return "", fmt.Errorf("websites.Fetch: URL \"%s\" is unreachable.", url.String())
					}

				} else {
					return "", fmt.Errorf("websites.Fetch: Invalid URL \"%s\".", url.String())
				}

			} else {
				return "", fmt.Errorf("websites.Fetch: %s", err1.Error())
			}

		} else {
			return "", fmt.Errorf("websites.Fetch: %s", err0.Error())
		}

	} else {
		return "", fmt.Errorf("websites.Fetch: Invalid parameter \"format\" must be \"markdown\", \"text\" or \"html\".")
	}

}
