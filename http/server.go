package http

import (
	"context"
	stderrors "errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Server struct {
	address string
	srv     *http.Server
	router  chi.Router

	logger *zap.Logger
}

func NewServer(address string, logger *zap.Logger, opts ...ServerOption) (s *Server) {
	s = &Server{
		address: address,

		router: chi.NewRouter(),

		logger: logger,
	}

	// fork gin-zap and implement it for chi
	s.router.Use(middleware.Logger, middleware.RequestID, middleware.Recoverer)

	s.srv = &http.Server{
		Addr:              s.address,
		Handler:           s.router,
		TLSConfig:         nil,
		ReadTimeout:       DefaultReadTimeout,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
		WriteTimeout:      DefaultWriteTimeout,
		IdleTimeout:       DefaultIdleTimeout,
		MaxHeaderBytes:    DefaultMaxHeaderBytes,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	for _, opt := range opts {
		opt.apply(s)
	}

	return s
}

func (s Server) Address() string {
	return s.address
}

func (s *Server) ListenAndServe() (err error) {
	if err = s.srv.ListenAndServe(); !stderrors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "listen and server error")
	}

	return nil
}

func (s *Server) Close(ctx context.Context) (err error) {
	return s.srv.Shutdown(ctx)
}
