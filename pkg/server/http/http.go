package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server struct {
	*gin.Engine
	httpSrv *http.Server
	host    string
	port    int
	logger  *zerolog.Logger
}

type Option func(s *Server)

type RouterAttachable interface {
	AttachRouter(*gin.Engine)
}

func NewServer(engine *gin.Engine, logger *zerolog.Logger, opts ...Option) *Server {
	s := &Server{
		Engine: engine,
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

func (s *Server) AttachRouter(router *gin.Engine) {
	s.Engine = router
}

func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	s.httpSrv = &http.Server{
		Addr:              addr,
		Handler:           s.Engine,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	s.logger.Info().
		Str("host", s.host).
		Int("port", s.port).
		Msg("Starting HTTP server")

	go func() {
		if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Warn().Msg("Shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.httpSrv.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Server forced to shutdown")
		return err
	}
	return nil
}
