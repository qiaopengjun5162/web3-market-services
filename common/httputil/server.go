package httputil

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
)

type HTTPServer struct {
	listener net.Listener
	server   *http.Server
	closed   atomic.Bool
}

type HTTPOption func(s *HTTPServer) error

//
//func NewHTTPServer(addr string, handler http.Handler, opts ...HTTPOption) (*HTTPServer, error) {
//	listener, err := net.Listen("tcp", addr)
//	if err != nil {
//		log.Error("Failed to listen http server", "err", err)
//		return nil, errors.New("Failed to listen http server")
//	}
//
//	serverCtx, serverCancel := context.WithCancel(context.Background())
//	server := &http.Server{
//		Handler:           handler,
//		ReadTimeout:       timeouts.ReadTimeout,
//		WriteTimeout:      timeouts.WriteTimeout,
//		IdleTimeout:       timeouts.IdleTimeout,
//		ReadHeaderTimeout: timeouts.ReadHeaderTimeout,
//		BaseContext: func(listener net.Listener) context.Context {
//			return serverCtx
//		},
//	}
//	srv := &HTTPServer{
//		listener: listener,
//		server:   server,
//	}
//
//	for _, opt := range opts {
//		if err := opt(srv); err != nil {
//			serverCancel()
//			log.Error("Failed to apply option", "err", err)
//			return nil, errors.New("Failed to apply option")
//		}
//	}
//	go func() {
//		err := srv.server.Serve(listener)
//		serverCancel()
//		if errors.Is(err, http.ErrServerClosed) {
//			log.Info("HTTP server closed")
//			srv.closed.Store(true)
//		} else {
//			log.Error("Failed to start http server", "err", err)
//			panic("Failed to start http server")
//		}
//	}()
//	return srv, nil
//}

func NewHTTPServer(addr string, handler http.Handler, opts ...HTTPOption) (*HTTPServer, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("listen errorr=", err)
		return nil, errors.New("Init listener fail")
	}
	srvCtx, srvCancel := context.WithCancel(context.Background())
	srv := &http.Server{
		Handler:           handler,
		ReadTimeout:       timeouts.ReadTimeout,
		ReadHeaderTimeout: timeouts.ReadHeaderTimeout,
		WriteTimeout:      timeouts.WriteTimeout,
		IdleTimeout:       timeouts.IdleTimeout,
		BaseContext: func(listener net.Listener) context.Context {
			return srvCtx
		},
	}
	out := &HTTPServer{listener: listener, server: srv}

	for _, opt := range opts {
		if err := opt(out); err != nil {
			srvCancel()
			fmt.Println("apply err:", err)
			return nil, errors.New("One of http op fail")
		}
	}
	go func() {
		err := out.server.Serve(listener)
		srvCancel()
		if errors.Is(err, http.ErrServerClosed) {
			out.closed.Store(true)
		} else {
			fmt.Println("unknow err:", err)
			panic("unknow error")
		}
	}()
	return out, nil
}

func (s *HTTPServer) Addr() net.Addr {
	return s.listener.Addr()
}

func (s *HTTPServer) IsClosed() bool {
	return s.closed.Load()
}

func (s *HTTPServer) Wait() {
	<-s.server.BaseContext(s.listener).Done()
}

func (s *HTTPServer) Context() context.Context {
	return s.server.BaseContext(s.listener)
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	if err := s.Shutdown(ctx); err != nil {
		if errors.Is(err, ctx.Err()) {
			return s.Close()
		}
		return err
	}
	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	if s.closed.Load() {
		return nil
	}
	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) Close() error {
	return s.server.Close()
}

func WithMaxHeaderBytes(maxBytes int) HTTPOption {
	return func(s *HTTPServer) error {
		s.server.MaxHeaderBytes = maxBytes
		return nil
	}
}
