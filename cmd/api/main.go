package main

import (
	"context"
	"github.com/kaffein/goffy/pkg/app"
	"github.com/kaffein/goffy/pkg/config"
	"github.com/kaffein/goffy/pkg/logger"
	grpcServer "github.com/kaffein/goffy/pkg/server/grpc"
	httpServer "github.com/kaffein/goffy/pkg/server/http"

	"github.com/gin-gonic/gin"
)

func init() {
	config.NewConfig("./config/local.yml")
}

func main() {
	// Logger config
	log := logger.NewLogger()

	// Init Gin router
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// HTTP server
	httpSrv := httpServer.NewServer(router, log,
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
		app.WithServer(httpSrv, grpcSrv),
	)

	// Run app
	if err := application.Run(context.Background()); err != nil {
		log.Fatal().Msg("App exited with error")
	}
}
