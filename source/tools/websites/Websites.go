package websites

import "exocomp/schemas"
import utils_fmt "exocomp/utils/fmt"
import utils_capabilities "exocomp/utils/capabilities"
import "context"
import "fmt"
import net_http "net/http"
import "net"
import "slices"
import "time"

type Websites struct {
	Methods      []string
	Playground   string
	Sandbox      string
	Capabilities *utils_capabilities.Set
	client       *net_http.Client
}

func NewWebsites(methods []string, playground string, sandbox string) *Websites {

	return NewWebsitesWithCapabilities(methods, playground, sandbox, nil)

}

func NewWebsitesWithCapabilities(
	methods []string,
	playground string,
	sandbox string,
	capabilities *utils_capabilities.Set,
) *Websites {

	websites := &Websites{
		Methods:      methods,
		Playground:   playground,
		Sandbox:      sandbox,
		Capabilities: capabilities,
	}

	dialer := &net.Dialer{
		Timeout: 30 * time.Second,
	}

	websites.client = &net_http.Client{
		Timeout: 30 * time.Second,
		Transport: &net_http.Transport{
			Proxy: net_http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, network string, address string) (net.Conn, error) {

				host, _, err0 := net.SplitHostPort(address)

				if err0 != nil {
					host = address
				}

				if websites.canConnect(host) == false {
					return nil, fmt.Errorf("websites.Fetch: Missing network.Connect permission.")
				} else {
					return dialer.DialContext(ctx, network, address)
				}

			},
		},
		CheckRedirect: func(request *net_http.Request, via []*net_http.Request) error {

			if websites.canConnect(request.URL.Hostname()) == true {

				if len(via) >= 10 {
					return fmt.Errorf("websites.Fetch: Stopped after 10 redirects.")
				}

				return nil

			} else {
				return fmt.Errorf("websites.Fetch: Missing network.Connect permission.")
			}
		},
	}

	return websites

}

func (tool *Websites) Name() string {
	return "websites"
}

func (tool *Websites) Call(method string, arguments map[string]any) (string, error) {

	if tool.HasMethod(method) == true {

		if method == "Fetch" {

			url, ok1 := arguments["url"].(string)
			user_agent, ok2 := arguments["user_agent"].(string)
			format, ok3 := arguments["format"].(string)

			if ok1 == true && ok2 == true && ok3 == true {
				return tool.Fetch(url, user_agent, utils_fmt.FormatSingleLine(format))
			} else if ok1 == true && ok2 == true && ok3 == false {
				return tool.Fetch(url, user_agent, "markdown")
			} else if ok1 == true && ok2 == false && ok3 == true {
				return tool.Fetch(url, "", utils_fmt.FormatSingleLine(format))
			} else if ok1 == true && ok2 == false && ok3 == false {
				return tool.Fetch(url, "", "markdown")
			} else if ok1 == false && ok2 == true && ok3 == true {
				return "", fmt.Errorf("websites.%s: %s", method, "Invalid parameter \"url\" is not a string.")
			} else if ok1 == false && ok2 == true && ok3 == false {
				return "", fmt.Errorf("websites.%s: %s", method, "Invalid parameter \"url\" is not a string.")
			} else if ok1 == false && ok2 == false && ok3 == true {
				return "", fmt.Errorf("websites.%s: %s", method, "Invalid parameter \"url\" is not a string.")
			} else {
				return "", fmt.Errorf("websites.%s: %s", method, "Invalid parameter \"url\" is not a string.")
			}

		} else if method == "List" {

			return tool.List()

		} else if method == "Stat" {

			url, ok1 := arguments["url"].(string)
			user_agent, ok2 := arguments["user_agent"].(string)

			if ok1 == true && ok2 == true {
				return tool.Stat(url, user_agent)
			} else if ok1 == true && ok2 == false {
				return tool.Stat(url, "")
			} else if ok1 == false && ok2 == true {
				return "", fmt.Errorf("websites.%s: %s", method, "Invalid parameter \"url\" is not a string.")
			} else {
				return "", fmt.Errorf("websites.%s: %s", method, "Invalid parameter \"url\" is not a string.")
			}

		} else {
			return "", fmt.Errorf("websites.%s: Invalid method.", method)
		}

	} else {
		return "", fmt.Errorf("websites.%s: Method not allowed.", method)
	}

}

func (tool *Websites) GetContent(id string) (any, error) {
	return nil, nil
}

func (tool *Websites) GetContentIdentifiers() []string {
	return []string{}
}

func (tool *Websites) HasMethod(method string) bool {
	return slices.Contains(tool.Methods, method) == true
}

func (tool *Websites) Schemas() []schemas.Tool {

	result := make([]schemas.Tool, 0)

	for _, method := range tool.Methods {

		for _, schema := range WebsitesSchema {

			if schema.Function.Name == fmt.Sprintf("%s.%s", tool.Name(), method) {
				result = append(result, schema)
			}

		}

	}

	return result

}
