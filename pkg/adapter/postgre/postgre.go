package postgre

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

var (
	db   *sql.DB
	once sync.Once
)

func GetDB() *sql.DB {
	return db
}

type Adapter struct {
	host     string
	port     int
	username string
	password string
	dbname   string
	sslmode  string
	logger   *zerolog.Logger
}

type Option func(s *Adapter)

func NewDatabase(logger *zerolog.Logger, opts ...Option) *Adapter {
	s := &Adapter{
		logger: logger,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithHost(host string) Option {
	return func(s *Adapter) {
		s.host = host
	}
}

func WithPort(port int) Option {
	return func(s *Adapter) {
		s.port = port
	}
}

func WithUsername(username string) Option {
	return func(s *Adapter) {
		s.username = username
	}
}

func WithPassword(password string) Option {
	return func(s *Adapter) {
		s.password = password
	}
}

func WithDBName(dbname string) Option {
	return func(s *Adapter) {
		s.dbname = dbname
	}
}

func WithSSLMode(mode string) Option {
	return func(s *Adapter) {
		s.sslmode = mode
	}
}

// Fungsi untuk membuka koneksi PostgreSQL secara global
func (s *Adapter) Start(ctx context.Context) error {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			s.host, s.port, s.username, s.password, s.dbname, s.sslmode,
		)

		// Membuka koneksi ke PostgreSQL
		var err error
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			s.logger.Fatal().Err(err).Msg("Failed to open PostgreSQL connection")
			return
		}

		// Verifikasi koneksi dengan ping
		if err := db.PingContext(ctx); err != nil {
			s.logger.Fatal().Err(err).Msg("Failed to ping PostgreSQL database")
			return
		}

		s.logger.Info().Str("host", s.host).Int("port", s.port).Msg("Successfully connected to PostgreSQL")
	})
	return nil
}

// Fungsi untuk menghentikan koneksi
func (s *Adapter) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	done := make(chan error)
	go func() {
		done <- db.Close()
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

	log.Warn().Msg("PostgreSQL connection closed successfully")
	return nil
}
