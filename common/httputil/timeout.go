package httputil

import "time"

var DefaultTimeOuts = HTTPTimeouts{
	ReadTimeout:       time.Second * 15,
	ReadHeaderTimeout: time.Second * 15,
	WriteTimeout:      time.Second * 15,
	IdleTimeout:       time.Second * 120,
}

type HTTPTimeouts struct {
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func WithTimeouts(timeouts HTTPTimeouts) HTTPOption {
	return func(s *HTTPServer) error {
		s.server.ReadTimeout = timeouts.ReadTimeout
		s.server.ReadHeaderTimeout = timeouts.ReadHeaderTimeout
		s.server.WriteTimeout = timeouts.WriteTimeout
		s.server.IdleTimeout = timeouts.IdleTimeout
		return nil
	}
}
