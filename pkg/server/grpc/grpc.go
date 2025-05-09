package grpc

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Server struct {
	*grpc.Server
	host   string
	port   int
	logger *zerolog.Logger
}

type Option func(s *Server)

func NewServer(logger *zerolog.Logger, opts ...Option) *Server {
	s := &Server{
		Server: grpc.NewServer(),
		logger: logger,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
func WithServerHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}
func WithServerPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		s.logger.Fatal().Err(err).Msg("failed to listen")
	}

	s.logger.Info().
		Str("host", s.host).
		Int("port", s.port).
		Msg("Starting Grpc server")

	if err = s.Server.Serve(lis); err != nil {
		s.logger.Fatal().Err(err).Msg("failed to serve")
	}

	return nil

}
func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	s.Server.GracefulStop()

	s.logger.Info().Msg("Shutting down Grpc server...")

	return nil
}
