package logger

import (
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

	timeFormat := config.Conf.Env.TimeFormat
	out := os.Stdout

	if config.Conf.Env.Mode != "local" {
		file, err := os.OpenFile(config.Conf.Logger.LogFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		out = file
	}

	writer := zerolog.ConsoleWriter{
		Out:        out,
		TimeFormat: timeFormat,
		FormatTimestamp: func(i interface{}) string {
			t, ok := i.(time.Time)
			if !ok {
				return "invalid-time"
			}
			return t.In(loc).Format(timeFormat)
		},
	}

	logger := zerolog.New(writer).
		With().
		Timestamp().
		Logger()
	log.Logger = logger
	return &logger
}
