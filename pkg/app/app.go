package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kaffein/goffy/pkg/server"
	"github.com/rs/zerolog/log"
)

type App struct {
	name    string
	servers []server.Server
}

type Option func(a *App)

func NewApp(opts ...Option) *App {
	a := &App{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func WithServer(servers ...server.Server) Option {
	return func(a *App) {
		a.servers = append(a.servers, servers...)
	}
}

func WithName(name string) Option {
	return func(a *App) {
		a.name = name
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msgf("Starting app: %s", a.name)

	// Start all servers
	for _, srv := range a.servers {
		go func(srv server.Server) {
			if err := srv.Start(ctx); err != nil {
				log.Error().Err(err).Msg("Server start error")
				cancel() // Cancel context on failure
			}
		}(srv)
	}

	select {
	case sig := <-signals:
		log.Info().Str("signal", sig.String()).Msg("Received termination signal")
	case <-ctx.Done():
		log.Info().Msg("Context canceled, shutting down servers")
	}

	// Graceful stop
	for _, srv := range a.servers {
		if err := srv.Stop(ctx); err != nil {
			log.Error().Err(err).Msg("Server stop error")
		}
	}

	log.Warn().Msgf("App %s exited", a.name)
	return nil
}
