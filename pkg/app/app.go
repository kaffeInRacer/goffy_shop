package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kaffein/goffy/pkg/server"
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

	log.Printf("Starting app: %s", a.name)

	// Start all servers
	for _, srv := range a.servers {
		go func(srv server.Server) {
			if err := srv.Start(ctx); err != nil {
				log.Printf("Server start error: %v", err)
				cancel() // Cancel context on failure
			}
		}(srv)
	}

	select {
	case sig := <-signals:
		log.Printf("Received termination signal: %s", sig)
	case <-ctx.Done():
		log.Println("Context canceled, shutting down servers")
	}

	// Graceful stop
	for _, srv := range a.servers {
		if err := srv.Stop(ctx); err != nil {
			log.Printf("Server stop error: %v", err)
		}
	}

	log.Printf("App %s exited", a.name)
	return nil
}
