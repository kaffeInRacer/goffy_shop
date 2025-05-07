package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kaffein/goffy/pkg/adapter" // Interface Adapter
	"github.com/kaffein/goffy/pkg/route"
	"github.com/kaffein/goffy/pkg/server"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	name     string
	adapters []adapter.Adapter
	servers  []server.Server
	router   *gin.Engine
}

type Option func(a *App)

func NewApp(opts ...Option) *App {
	a := &App{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func WithAdapter(adapters ...adapter.Adapter) Option {
	return func(a *App) {
		a.adapters = append(a.adapters, adapters...)
	}
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

func WithRouter(router *gin.Engine) Option {
	return func(a *App) {
		a.router = router
	}
}

func (a *App) setupRoutes() {
	route.SetupRoutes(a.router)
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msgf("Starting app: %s", a.name)

	for _, adapterInst := range a.adapters {
		if err := adapterInst.Start(ctx); err != nil {
			log.Err(err).Msg("Adapter start error")
			cancel()
			return err
		}
	}

	var wg sync.WaitGroup
	for _, srv := range a.servers {
		wg.Add(1)
		go func(s server.Server) {
			defer wg.Done()
			if err := srv.Start(ctx); err != nil {
				log.Err(err).Msg("Server start error")
				cancel()
			}
		}(srv)
	}
	wg.Wait()

	// ⚠ Recommendation: call a.setupRoutes() BEFORE the loop that launches server.Start()
	// to ensure all routes are ready before the server begins handling traffic.
	// Running it after server startup risks race conditions where routes might not
	// be fully registered before requests arrive.
	//
	// However, if this approach works reliably in your setup (e.g., due to internal
	// startup delays or no early traffic), you may ignore this warning.
	a.setupRoutes()

	select {
	case sig := <-signals:
		log.Info().Msgf("Received termination signal: %s", sig)
	case <-ctx.Done():
		log.Info().Msg("Context canceled, shutting down")
	}

	for _, srv := range a.servers {
		if err := srv.Stop(ctx); err != nil {
			log.Err(err).Msg("Server stop error")
		}
	}

	for _, adapterInst := range a.adapters {
		if err := adapterInst.Stop(ctx); err != nil {
			log.Err(err).Msg("Adapter stop error")
		}
	}

	wg.Wait()
	log.Warn().Msgf("App %s exited", a.name)
	return nil
}
