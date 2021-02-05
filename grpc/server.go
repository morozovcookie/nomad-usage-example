package grpc

import (
	"net"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	DefaultNetwork = "tcp"
)

type Server struct {
	address string
	srv     *grpc.Server

	logger *zap.Logger
}

func NewServer(address string, logger *zap.Logger) *Server {
	return &Server{
		address: address,
		srv:     grpc.NewServer(),

		logger: logger,
	}
}

func (s Server) Address() string {
	return s.address
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen(DefaultNetwork, s.address)
	if err != nil {
		return errors.Wrap(err, "listen error")
	}

	if err := s.srv.Serve(ln); err != nil {
		return errors.Wrap(err, "serve error")
	}

	return nil
}

func (s *Server) Close() {
	s.srv.GracefulStop()
}

func (s *Server) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.srv.RegisterService(desc, impl)
}
