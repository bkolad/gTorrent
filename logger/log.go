package logger

import (
	"os"

	z "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(z.ConsoleWriter{Out: os.Stdout})
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Debug(msg string) {
	log.Debug().Msg(msg)
}

func Error(msg string) {
	log.Error().Msg(msg)
}
