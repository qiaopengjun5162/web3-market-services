package httputil

import "net/http"

var timeouts = DefaultTimeOuts

func NewHttpServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler:           handler,
		ReadTimeout:       timeouts.ReadTimeout,
		WriteTimeout:      timeouts.WriteTimeout,
		IdleTimeout:       timeouts.IdleTimeout,
		ReadHeaderTimeout: timeouts.ReadHeaderTimeout,
	}
}
