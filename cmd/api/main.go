// main.go

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kaffein/goffy/pkg/adapter/postgre"
	"github.com/kaffein/goffy/pkg/app"
	"github.com/kaffein/goffy/pkg/config"
	"github.com/kaffein/goffy/pkg/logger"
	grpcServer "github.com/kaffein/goffy/pkg/server/gRPC"
	httpServer "github.com/kaffein/goffy/pkg/server/net_http"
)

func init() {
	config.LoadConfig("./config/local.yml")
}

func main() {
	log := logger.NewLogger()

	// engine
	engine := gin.Default()

	// Setup PostgreSQL
	dbPg := postgre.NewDatabase(log,
		postgre.WithHost(config.Conf.Database.Postgres.Host),
		postgre.WithPort(config.Conf.Database.Postgres.Port),
		postgre.WithUsername(config.Conf.Database.Postgres.User),
		postgre.WithPassword(config.Conf.Database.Postgres.Password),
		postgre.WithDBName(config.Conf.Database.Postgres.Name),
		postgre.WithSSLMode(config.Conf.Database.Postgres.SSL),
	)

	// Setup HTTP server
	httpSrv := httpServer.NewServer(engine, log,
		httpServer.WithServerHost(config.Conf.Server.HTTP.Host),
		httpServer.WithServerPort(config.Conf.Server.HTTP.Port),
	)

	// Setup gRPC server
	grpcSrv := grpcServer.NewServer(log,
		grpcServer.WithServerHost(config.Conf.Server.GRPC.Host),
		grpcServer.WithServerPort(config.Conf.Server.GRPC.Port),
	)

	// Create application with all servers, adaptors, and router
	application := app.NewApp(
		app.WithName("goffy"),
		app.WithAdapter(dbPg),
		app.WithServer(httpSrv, grpcSrv),
		app.WithRouter(engine),
	)

	// Run the application
	if err := application.Run(context.Background()); err != nil {
		log.Fatal().Msg("App exited with error")
	}
}
