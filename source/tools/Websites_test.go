package tools

import "fmt"
import "net/http"
import "net/http/httptest"
import "strings"
import "testing"

func TestWebsites_List(t *testing.T) {

	tool := NewWebsites("/tmp", "/tmp")

	result, err0 := tool.List()

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if !strings.HasPrefix(result, "websites.List:") {
		t.Errorf("Expected a websites.List report, got %s", result)
	}

	if !strings.Contains(result, "chrome-windows") {
		t.Errorf("Expected chrome-windows preset to be listed")
	}

}

func TestWebsites_Fetch_Markdown(t *testing.T) {

	var got_user_agent string
	var got_sec_ch_ua  string

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		got_user_agent = request.Header.Get("User-Agent")
		got_sec_ch_ua  = request.Header.Get("Sec-CH-UA")

		response.Header().Set("Content-Type", "text/html")

		fmt.Fprint(response, "<html><head><title>Test</title></head><body><h1>Hello</h1><p>World</p></body></html>")

	}))
	defer server.Close()

	tool := NewWebsites("/tmp", "/tmp")

	result, err0 := tool.Fetch(server.URL, "chrome-windows", "markdown")

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if !strings.Contains(result, "# Hello") {
		t.Errorf("Expected markdown heading, got:\n%s", result)
	}

	if !strings.Contains(result, "World") {
		t.Errorf("Expected markdown content, got:\n%s", result)
	}

	if got_user_agent == "" {
		t.Errorf("Expected a User-Agent header to be sent")
	}

	if got_sec_ch_ua == "" {
		t.Errorf("Expected a Sec-CH-UA header to be sent")
	}

}

func TestWebsites_Fetch_Text(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "text/html")
		fmt.Fprint(response, "<html><body><h1>Heading</h1><p>Body text</p></body></html>")
	}))
	defer server.Close()

	tool := NewWebsites("/tmp", "/tmp")

	result, err0 := tool.Fetch(server.URL, "", "text")

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if !strings.Contains(result, "Heading") || !strings.Contains(result, "Body text") {
		t.Errorf("Expected plain text content, got:\n%s", result)
	}

	if strings.Contains(result, "# ") {
		t.Errorf("Expected plain text without markdown syntax, got:\n%s", result)
	}

}

func TestWebsites_Fetch_InvalidURL(t *testing.T) {

	tool := NewWebsites("/tmp", "/tmp")

	_, err0 := tool.Fetch("ftp://example.com/file", "", "markdown")

	if err0 == nil {
		t.Errorf("Expected a non-nil error for an invalid URL scheme")
	}

}

func TestWebsites_Fetch_InvalidUserAgent(t *testing.T) {

	tool := NewWebsites("/tmp", "/tmp")

	_, err0 := tool.Fetch("https://example.com/", "does-not-exist", "markdown")

	if err0 == nil {
		t.Errorf("Expected a non-nil error for an unknown User-Agent")
	}

}

func TestWebsites_Stat(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("X-Custom", "yes")
		response.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tool := NewWebsites("/tmp", "/tmp")

	result, err0 := tool.Stat(server.URL, "")

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if !strings.Contains(result, "200 OK") {
		t.Errorf("Expected a 200 status, got:\n%s", result)
	}

	if !strings.Contains(result, "X-Custom") {
		t.Errorf("Expected the custom header to be reported, got:\n%s", result)
	}

}

func TestWebsites_Get(t *testing.T) {

	tool := NewWebsites("/tmp", "/tmp")

	content, err0 := tool.Get("chrome-windows")

	if err0 != nil {
		t.Errorf("Expected %v to be nil", err0)
	}

	if content == nil {
		t.Errorf("Expected a User-Agent, got nil")
	}

}
