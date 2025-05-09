package logger

import (
	"fmt"
	"github.com/kaffein/goffy/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func NewLogger() *zerolog.Logger {
	loc, err := time.LoadLocation(config.Conf.Env.TimeZone)
	if err != nil {
		panic(err)
	}
	
	out := os.Stdout

	if config.Conf.Env.Mode != "local" {
		file, err := os.OpenFile(config.Conf.Logger.LogFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		out = file
	}

	writer := zerolog.ConsoleWriter{
		Out:          out,
		TimeFormat:   config.Conf.Env.TimeFormat,
		TimeLocation: loc,
		FormatTimestamp: func(i interface{}) string {
			return fmt.Sprintf("[ %v ]", i)
		},
	}

	logger := zerolog.New(writer).
		With().
		Timestamp().
		Logger()
	log.Logger = logger
	return &logger
}
