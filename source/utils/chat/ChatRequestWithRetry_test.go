package chat

import "context"
import "errors"
import "io"
import net_http "net/http"
import "strings"
import "testing"

type requestMockResponse struct {
	status int
	header net_http.Header
	body   string
	err    error
}

type requestMockTransport struct {
	responses []requestMockResponse
	calls     int
	bodies    []string
	headers   []net_http.Header
}

func (transport *requestMockTransport) RoundTrip(request *net_http.Request) (*net_http.Response, error) {

	transport.calls++

	payload, _ := io.ReadAll(request.Body)
	transport.bodies = append(transport.bodies, string(payload))

	header := make(net_http.Header)

	for key, values := range request.Header {
		header[key] = values
	}

	transport.headers = append(transport.headers, header)

	index := transport.calls - 1

	if index >= len(transport.responses) {
		index = len(transport.responses) - 1
	}

	item := transport.responses[index]

	if item.err != nil {

		return nil, item.err

	} else {

		response_header := make(net_http.Header)

		for key, values := range item.header {
			response_header[key] = values
		}

		return &net_http.Response{
			StatusCode: item.status,
			Header:     response_header,
			Body:       io.NopCloser(strings.NewReader(item.body)),
			Request:    request,
		}, nil

	}

}

func TestChatRequestWithRetry_RetriesTransientStatus(t *testing.T) {

	transport := &requestMockTransport{
		responses: []requestMockResponse{
			{
				status: net_http.StatusServiceUnavailable,
				header: net_http.Header{"Retry-After": []string{"0"}},
				body:   "unavailable",
			},
			{
				status: net_http.StatusBadGateway,
				header: net_http.Header{"Retry-After": []string{"0"}},
				body:   "bad gateway",
			},
			{
				status: net_http.StatusOK,
				header: net_http.Header{},
				body:   "ok",
			},
		},
	}

	client := &net_http.Client{Transport: transport}
	payload := []byte(`{"model":"test-model"}`)

	response, err := ChatRequestWithRetry(client, payload, "http://localhost/v1/chat/completions", "token")

	if err == nil {

		if response.StatusCode != net_http.StatusOK {
			t.Errorf("Expected status %d to be %d", response.StatusCode, net_http.StatusOK)
		}

		if transport.calls != 3 {
			t.Errorf("Expected %d calls to be %d", transport.calls, 3)
		}

		for _, body := range transport.bodies {

			if body != string(payload) {
				t.Errorf("Expected request body %q to be %q", body, string(payload))
			}

		}

		if len(transport.headers) > 0 {

			if transport.headers[0].Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type %q to be %q", transport.headers[0].Get("Content-Type"), "application/json")
			}

			if transport.headers[0].Get("Accept") != "application/json" {
				t.Errorf("Expected Accept %q to be %q", transport.headers[0].Get("Accept"), "application/json")
			}

			if transport.headers[0].Get("Authorization") != "Bearer token" {
				t.Errorf("Expected Authorization %q to be %q", transport.headers[0].Get("Authorization"), "Bearer token")
			}

		} else {
			t.Errorf("Expected request headers to be recorded")
		}

		response_payload, err1 := io.ReadAll(response.Body)
		response.Body.Close()

		if err1 == nil {

			if string(response_payload) != "ok" {
				t.Errorf("Expected response body %q to be %q", string(response_payload), "ok")
			}

		} else {
			t.Errorf("Expected %v to be nil", err1)
		}

	} else {
		t.Errorf("Expected %v to be nil", err)
	}

}

func TestChatRequestWithRetry_DoesNotRetryClientError(t *testing.T) {

	transport := &requestMockTransport{
		responses: []requestMockResponse{
			{
				status: net_http.StatusBadRequest,
				header: net_http.Header{},
				body:   "bad request",
			},
		},
	}

	client := &net_http.Client{Transport: transport}

	response, err := ChatRequestWithRetry(client, []byte(`{}`), "http://localhost/v1/chat/completions", "")

	if err == nil {

		if response.StatusCode != net_http.StatusBadRequest {
			t.Errorf("Expected status %d to be %d", response.StatusCode, net_http.StatusBadRequest)
		}

		if transport.calls != 1 {
			t.Errorf("Expected %d calls to be %d", transport.calls, 1)
		}

		if len(transport.headers) > 0 {

			if transport.headers[0].Get("Authorization") != "" {
				t.Errorf("Expected no Authorization header, got %q", transport.headers[0].Get("Authorization"))
			}

		}

		response_payload, err1 := io.ReadAll(response.Body)
		response.Body.Close()

		if err1 == nil {

			if string(response_payload) != "bad request" {
				t.Errorf("Expected response body %q to be %q", string(response_payload), "bad request")
			}

		} else {
			t.Errorf("Expected %v to be nil", err1)
		}

	} else {
		t.Errorf("Expected %v to be nil", err)
	}

}

func TestChatRequestWithRetry_StopsAfterMaxAttempts(t *testing.T) {

	transport := &requestMockTransport{
		responses: []requestMockResponse{
			{
				status: net_http.StatusServiceUnavailable,
				header: net_http.Header{"Retry-After": []string{"0"}},
				body:   "unavailable",
			},
		},
	}

	client := &net_http.Client{Transport: transport}

	response, err := ChatRequestWithRetry(client, []byte(`{}`), "http://localhost/v1/chat/completions", "")

	if err == nil {

		if response.StatusCode != net_http.StatusServiceUnavailable {
			t.Errorf("Expected status %d to be %d", response.StatusCode, net_http.StatusServiceUnavailable)
		}

		if transport.calls != RetryMaxAttempts {
			t.Errorf("Expected %d calls to be %d", transport.calls, RetryMaxAttempts)
		}

		response_payload, err1 := io.ReadAll(response.Body)
		response.Body.Close()

		if err1 == nil {

			if string(response_payload) != "unavailable" {
				t.Errorf("Expected response body %q to be %q", string(response_payload), "unavailable")
			}

		} else {
			t.Errorf("Expected %v to be nil", err1)
		}

	} else {
		t.Errorf("Expected %v to be nil", err)
	}

}

func TestChatRequestWithRetry_RetriesTransportError(t *testing.T) {

	transport := &requestMockTransport{
		responses: []requestMockResponse{
			{
				err: errors.New("connection reset"),
			},
			{
				status: net_http.StatusOK,
				header: net_http.Header{},
				body:   "ok",
			},
		},
	}

	client := &net_http.Client{Transport: transport}

	response, err := ChatRequestWithRetry(client, []byte(`{}`), "http://localhost/v1/chat/completions", "")

	if err == nil {

		if response.StatusCode != net_http.StatusOK {
			t.Errorf("Expected status %d to be %d", response.StatusCode, net_http.StatusOK)
		}

		if transport.calls != 2 {
			t.Errorf("Expected %d calls to be %d", transport.calls, 2)
		}

		response.Body.Close()

	} else {
		t.Errorf("Expected %v to be nil", err)
	}

}

func TestChatRequestWithRetry_DoesNotRetryCancelledContext(t *testing.T) {

	transport := &requestMockTransport{
		responses: []requestMockResponse{
			{
				err: context.Canceled,
			},
		},
	}

	client := &net_http.Client{Transport: transport}

	response, err := ChatRequestWithRetry(client, []byte(`{}`), "http://localhost/v1/chat/completions", "")

	if err != nil {

		if errors.Is(err, context.Canceled) == false {
			t.Errorf("Expected a cancelled context error, got %v", err)
		}

		if response != nil {
			t.Errorf("Expected %v to be nil", response)
		}

	} else {
		t.Errorf("Expected an error")
	}

	if transport.calls != 1 {
		t.Errorf("Expected %d calls to be %d", transport.calls, 1)
	}

}
