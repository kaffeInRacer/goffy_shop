package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kaffein/goffy/pkg/adapter"
	"github.com/kaffein/goffy/pkg/route"
	"github.com/kaffein/goffy/pkg/server"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	Name     string
	Adapters []adapter.Adapter
	Servers  []server.Server
	QuitChan chan os.Signal
	Router   *gin.Engine
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
		a.Adapters = append(a.Adapters, adapters...)
	}
}

func WithServer(servers ...server.Server) Option {
	return func(a *App) {
		a.Servers = append(a.Servers, servers...)

func WithName(name string) Option {
	return func(a *App) {
		a.Name = name
	}
}

func WithRouter(router *gin.Engine) Option {
	return func(a *App) {
		a.Router = router
	}
}

func (a *App) setupRoutes() {
	route.SetupRoutes(a.Router)
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if a.QuitChan == nil {
		a.QuitChan = make(chan os.Signal, 1)
	}
	signal.Notify(a.QuitChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msgf("Starting app: %s", a.Name)

	for _, adapterInst := range a.Adapters {
		if err := adapterInst.Start(ctx); err != nil {
			log.Err(err).Msg("Adapter start error")
			cancel()
			return err
		}
	}

	var wg sync.WaitGroup
	for _, srv := range a.Servers {
		wg.Add(1)
		go func(s server.Server) {
			defer wg.Done()
			if err := s.Start(ctx); err != nil {
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
	if a.Router != nil {
		a.setupRoutes()
	}

	select {
	case sig := <-a.QuitChan:
		log.Info().Msgf("Received termination signal: %s", sig)
	case <-ctx.Done():
		log.Info().Msg("Context canceled, shutting down")
	}

	for _, srv := range a.Servers {
		if err := srv.Stop(ctx); err != nil {
			log.Err(err).Msg("Server stop error")
		}
	}

	for _, adapterInst := range a.Adapters {
		if err := adapterInst.Stop(ctx); err != nil {
			log.Err(err).Msg("Adapter stop error")
		}
	}

	wg.Wait()
	log.Warn().Msgf("App %s exited", a.Name)
	return nil
}
