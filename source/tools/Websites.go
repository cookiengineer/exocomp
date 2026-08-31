package tools

import "exocomp/schemas"
import utils_fmt "exocomp/utils/fmt"
import "exocomp/types"
import "fmt"
import "io"
import net_http "net/http"
import "slices"
import "sort"
import "strings"
import "time"

type Websites struct {
	Methods    []string
	Playground string
	Sandbox    string
	client     *net_http.Client
}

func NewWebsites(methods []string, playground string, sandbox string) *Websites {

	websites := &Websites{
		Methods:    methods,
		Playground: playground,
		Sandbox:    sandbox,
		client:     &net_http.Client{
			Timeout: 30 * time.Second,
		},
	}

	return websites

}

func (tool *Websites) Name() string {
	return "websites"
}

func (tool *Websites) Call(method string, arguments map[string]interface{}) (string, error) {

	if slices.Contains(tool.Methods, method) == true {

		if method == "Fetch" {

			url,        ok1 := arguments["url"].(string)
			user_agent, ok2 := arguments["user_agent"].(string)
			format,     ok3 := arguments["format"].(string)

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

			url,        ok1 := arguments["url"].(string)
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

func (tool *Websites) Fetch(url_str string, user_agent string, format string) (string, error) {

	format = strings.ToLower(strings.TrimSpace(format))

	if format == "" {
		format = "markdown"
	} else if format != "markdown" && format != "text" && format != "html" {
		return "", fmt.Errorf("websites.Fetch: Invalid parameter \"format\" must be \"markdown\", \"text\" or \"html\".")
	}

	url, err0 := parseWebsiteURL(url_str)

	if err0 != nil {
		return "", err0
	}

	useragent, err1 := getWebsiteUserAgent(user_agent)

	if err1 != nil {
		return "", err1
	}

	request, err2 := net_http.NewRequest(net_http.MethodGet, url.String(), nil)

	if err2 != nil {
		return "", fmt.Errorf("websites.Fetch: %s", err2.Error())
	}

	request.Header = useragent.Header()

	accept          := request.Header.Get("Accept")
	accept_language := request.Header.Get("Accept-Language")

	if accept == "" {
		accept = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
	}

	if accept_language == "" {
		accept_language = "en-US,en;q=0.9"
	}

	request.Header.Set("Accept",          accept)
	request.Header.Set("Accept-Language", accept_language)

	response, err3 := tool.client.Do(request)

	if err3 != nil {
		return "", fmt.Errorf("websites.Fetch: %s", err3.Error())
	}

	defer response.Body.Close()

	bytes, err4 := io.ReadAll(io.LimitReader(response.Body, 16*1024*1024))

	if err4 != nil {
		return "", fmt.Errorf("websites.Fetch: %s", err4.Error())
	}

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
		fmt.Sprintf("websites.Fetch: %s (%s)", url.String(), response.Status),
		strings.TrimSpace(content),
	}, "\n"), nil

}

func (tool *Websites) GetContent(id string) (any, error) {
	return nil, nil
}

func (tool *Websites) List() (string, error) {

	names := make([]string, 0)

	for _, useragent := range types.UserAgents {
		names = append(names, useragent.Name)
	}

	sort.Strings(names)

	lines := make([]string, 0)
	lines = append(lines, fmt.Sprintf("websites.List: %d User-Agent presets available.", len(names)))

	for _, name := range names {

		useragent := types.GetUserAgent(name)

		if useragent != nil {
			lines = append(lines, fmt.Sprintf("- Name: \"%s\", Platform: %s, Mobile: %t", useragent.Name, useragent.Platform, useragent.Mobile))
		}

	}

	return strings.Join(lines, "\n"), nil

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

func (tool *Websites) Stat(url_str string, user_agent string) (string, error) {

	url, err0 := parseWebsiteURL(url_str)

	if err0 != nil {
		return "", err0
	}

	useragent, err1 := getWebsiteUserAgent(user_agent)

	if err1 != nil {
		return "", err1
	}

	request, err2 := net_http.NewRequest(net_http.MethodHead, url.String(), nil)

	if err2 != nil {
		return "", fmt.Errorf("websites.Stat: %s", err2.Error())
	}

	request.Header = useragent.Header()

	response, err3 := tool.client.Do(request)

	if err3 != nil {

		// Some servers reject HEAD, fall back to a body-less GET
		request2, err4 := net_http.NewRequest(net_http.MethodGet, url.String(), nil)

		if err4 != nil {
			return "", fmt.Errorf("websites.Stat: %s", err4.Error())
		}

		request2.Header = useragent.Header()

		response2, err5 := tool.client.Do(request2)

		if err5 != nil {
			return "", fmt.Errorf("websites.Stat: %s", err5.Error())
		}

		io.Copy(io.Discard, response2.Body)
		response2.Body.Close()

		response = response2

	} else {

		io.Copy(io.Discard, response.Body)
		response.Body.Close()

	}

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

}

