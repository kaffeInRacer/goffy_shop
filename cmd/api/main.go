package main

import (
	"context"
	"github.com/kaffein/goffy/pkg/app"
	"github.com/kaffein/goffy/pkg/config"
	"github.com/kaffein/goffy/pkg/db/postgres"
	"github.com/kaffein/goffy/pkg/logger"
	"github.com/kaffein/goffy/pkg/route"
	grpcServer "github.com/kaffein/goffy/pkg/server/gRPC"
	httpServer "github.com/kaffein/goffy/pkg/server/http"
)

func init() {
	config.LoadConfig("./config/local.yml")
}

func main() {
	// Logger config
	log := logger.NewLogger()

	// (db) PostgreSQL
	pgSrv := postgres.NewDatabase(log,
		postgres.WithHost(config.Conf.Database.Postgres.Host),
		postgres.WithPort(config.Conf.Database.Postgres.Port),
		postgres.WithUsername(config.Conf.Database.Postgres.User),
		postgres.WithPassword(config.Conf.Database.Postgres.Password),
		postgres.WithDBName(config.Conf.Database.Postgres.Name),
		postgres.WithSSLMode(config.Conf.Database.Postgres.SSL),
	)

	// HTTP server
	httpSrv := httpServer.NewServer(
		route.NewRouter(pgSrv.DB), log,
		httpServer.WithServerHost(config.Conf.Server.HTTP.Host),
		httpServer.WithServerPort(config.Conf.Server.HTTP.Port),
	)

	// gRPC server
	grpcSrv := grpcServer.NewServer(log,
		grpcServer.WithServerHost(config.Conf.Server.GRPC.Host),
		grpcServer.WithServerPort(config.Conf.Server.GRPC.Port),
	)

	// Compose app with servers
	application := app.NewApp(
		app.WithName("goffy"),
		app.WithServer(httpSrv, grpcSrv, pgSrv),
	)

	// Run app
	if err := application.Run(context.Background()); err != nil {
		log.Fatal().Msg("App exited with error")
	}
}
