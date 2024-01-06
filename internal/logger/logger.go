package logger

import (
	"errors"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Log zerolog.Logger = zerolog.New(os.Stdout)

type Level int

const (
	Debug Level = iota
	Info
)

var ErrNoSuchLevel = errors.New("No such logging level")

func SetLevel(lvl Level) error {
	switch lvl {
	case Debug:
		Log = log.Level(zerolog.DebugLevel)
	case Info:
		Log = log.Level(zerolog.InfoLevel)
	default:
		return ErrNoSuchLevel
	}

	return nil
}
