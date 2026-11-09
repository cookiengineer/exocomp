package http

import "net/http"
import "sync/atomic"

type handlerValue struct {
	handler http.Handler
}

type Handler struct {
	value atomic.Value // stores http.Handler
}

func NewHandler(http_handler http.Handler) *Handler {

	handler := &Handler{}
	handler.value.Store(&handlerValue{
		handler: http_handler,
	})

	return handler

}

func (handler *Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	value := handler.value.Load().(*handlerValue)
	value.handler.ServeHTTP(response, request)

}

func (handler *Handler) Set(http_handler http.Handler) {

	handler.value.Store(&handlerValue{
		handler: http_handler,
	})

}
