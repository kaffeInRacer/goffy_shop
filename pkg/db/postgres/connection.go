package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type Server struct {
	DB       *sql.DB
	host     string
	port     int
	username string
	password string
	dbname   string
	sslmode  string
	logger   *zerolog.Logger
}

type Option func(s *Server)

func NewDatabase(logger *zerolog.Logger, opts ...Option) *Server {
	s := &Server{
		logger: logger,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func WithUsername(username string) Option {
	return func(s *Server) {
		s.username = username
	}
}

func WithPassword(password string) Option {
	return func(s *Server) {
		s.password = password
	}
}

func WithDBName(dbname string) Option {
	return func(s *Server) {
		s.dbname = dbname
	}
}

func WithSSLMode(mode string) Option {
	return func(s *Server) {
		s.sslmode = mode
	}
}

func (s *Server) Start(ctx context.Context) error {
	var once sync.Once
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			s.host, s.port, s.username, s.password, s.dbname, s.sslmode,
		)

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			s.logger.Fatal().Err(err).Msg("Failed to open PostgreSQL connection")
		}

		if err := db.PingContext(ctx); err != nil {
			s.logger.Fatal().Err(err).Msg("Failed to ping PostgreSQL database")
		}

		s.DB = db
		s.logger.Info().
			Str("host", s.host).
			Int("port", s.port).
			Msg("Starting postgreSQL connection")
	})
	return nil
}

func (s *Server) Stop(ctx context.Context) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	done := make(chan error)
	go func() {
		done <- s.DB.Close()
	}()

	select {
	case <-ctx.Done():
		s.logger.Error().Msg("Timeout while closing PostgreSQL connection")
		return ctx.Err()
	case err := <-done:
		if err != nil {
			s.logger.Error().Err(err).Msg("Error closing PostgreSQL connection")
			return err
		}
	}

	log.Warn().Msg("Stopping postgreSQL connection...")
	return nil
}
