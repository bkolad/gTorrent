package logger

import (
	"os"

	z "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	//Warn()
	//Error()
}

var Default Logger

func init() {
	log.Logger = log.Output(z.ConsoleWriter{Out: os.Stdout})
	Default = zerolog{}
}

type zerolog struct {
}

func (z zerolog) Info(msg string) {
	log.Info().Msg(msg)
}

func (z zerolog) Debug(msg string) {
	log.Debug().Msg(msg)
}
