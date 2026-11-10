package engine

import utils_chat "exocomp/utils/chat"
import "context"
import "errors"
import "io"
import net_http "net/http"
import "strings"
import "testing"

type retryMockResponse struct {
	status int
	header net_http.Header
	body   string
	err    error
}

type retryMockTransport struct {
	responses []retryMockResponse
	calls     int
	bodies    []string
}

func (transport *retryMockTransport) RoundTrip(request *net_http.Request) (*net_http.Response, error) {

	transport.calls++

	payload, _ := io.ReadAll(request.Body)
	transport.bodies = append(transport.bodies, string(payload))

	index := transport.calls - 1

	if index >= len(transport.responses) {
		index = len(transport.responses) - 1
	}

	item := transport.responses[index]

	if item.err != nil {

		return nil, item.err

	} else {

		header := make(net_http.Header)

		for key, values := range item.header {
			header[key] = values
		}

		return &net_http.Response{
			StatusCode: item.status,
			Header:     header,
			Body:       io.NopCloser(strings.NewReader(item.body)),
			Request:    request,
		}, nil

	}

}

func newRetryTestSession(transport *retryMockTransport) *Session {

	session := newTestSession()
	session.client = &net_http.Client{Transport: transport}
	session.Console.DisableMethod("Warn")
	session.Console.DisableMethod("Error")

	return session

}

func TestSession_InferChatCompletions_RetriesTransientStatus(t *testing.T) {

	transport := &retryMockTransport{
		responses: []retryMockResponse{
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
				body:   `{"choices":[{"message":{"role":"assistant","content":"done"}}]}`,
			},
		},
	}

	session := newRetryTestSession(transport)
	err := session.inferChatCompletions()

	if err == nil {

		if transport.calls != 3 {
			t.Errorf("Expected %d calls to be %d", transport.calls, 3)
		}

		if len(transport.bodies) != transport.calls {
			t.Errorf("Expected %d recorded bodies to be %d", len(transport.bodies), transport.calls)
		}

		for _, body := range transport.bodies {

			if body != transport.bodies[0] {
				t.Errorf("Expected request body %q to be %q", body, transport.bodies[0])
			}

		}

		last := session.GetLastMessage()

		if last != nil {

			if last.Content != "done" {
				t.Errorf("Expected response content %q to be %q", last.Content, "done")
			}

		} else {
			t.Fatalf("Expected a last message to be appended")
		}

	} else {
		t.Errorf("Expected %v to be nil", err)
	}

}

func TestSession_InferChatCompletions_DoesNotRetryClientError(t *testing.T) {

	transport := &retryMockTransport{
		responses: []retryMockResponse{
			{
				status: net_http.StatusBadRequest,
				header: net_http.Header{},
				body:   `{"error":{"message":"bad request","type":"invalid_request_error"}}`,
			},
		},
	}

	session := newRetryTestSession(transport)
	err := session.inferChatCompletions()

	if err != nil {

		if strings.Contains(err.Error(), "incompatible with OpenAPI") == false {
			t.Errorf("Expected an OpenAPI incompatibility error, got %v", err)
		}

	} else {
		t.Errorf("Expected an error")
	}

	if transport.calls != 1 {
		t.Errorf("Expected %d calls to be %d", transport.calls, 1)
	}

}

func TestSession_InferChatCompletions_StopsAfterMaxAttempts(t *testing.T) {

	transport := &retryMockTransport{
		responses: []retryMockResponse{
			{
				status: net_http.StatusServiceUnavailable,
				header: net_http.Header{"Retry-After": []string{"0"}},
				body:   "unavailable",
			},
		},
	}

	session := newRetryTestSession(transport)
	err := session.inferChatCompletions()

	if err != nil {

		if strings.Contains(err.Error(), "unexpected HTTP error 503") == false {
			t.Errorf("Expected an unexpected HTTP error, got %v", err)
		}

	} else {
		t.Errorf("Expected an error")
	}

	if transport.calls != utils_chat.RetryMaxAttempts {
		t.Errorf("Expected %d calls to be %d", transport.calls, utils_chat.RetryMaxAttempts)
	}

}

func TestSession_InferChatCompletions_RetriesTransportError(t *testing.T) {

	transport := &retryMockTransport{
		responses: []retryMockResponse{
			{
				err: errors.New("connection reset"),
			},
			{
				status: net_http.StatusOK,
				header: net_http.Header{},
				body:   `{"choices":[{"message":{"role":"assistant","content":"done"}}]}`,
			},
		},
	}

	session := newRetryTestSession(transport)
	err := session.inferChatCompletions()

	if err == nil {

		if transport.calls != 2 {
			t.Errorf("Expected %d calls to be %d", transport.calls, 2)
		}

		last := session.GetLastMessage()

		if last != nil {

			if last.Content != "done" {
				t.Errorf("Expected response content %q to be %q", last.Content, "done")
			}

		} else {
			t.Fatalf("Expected a last message to be appended")
		}

	} else {
		t.Errorf("Expected %v to be nil", err)
	}

}

func TestSession_InferChatCompletions_DoesNotRetryCancelledContext(t *testing.T) {

	transport := &retryMockTransport{
		responses: []retryMockResponse{
			{
				err: context.Canceled,
			},
		},
	}

	session := newRetryTestSession(transport)
	err := session.inferChatCompletions()

	if err != nil {

		if errors.Is(err, context.Canceled) == false {
			t.Errorf("Expected a cancelled context error, got %v", err)
		}

	} else {
		t.Errorf("Expected an error")
	}

	if transport.calls != 1 {
		t.Errorf("Expected %d calls to be %d", transport.calls, 1)
	}

}
