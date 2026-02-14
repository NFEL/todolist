package logger

import (
	"log/slog"

	"github.com/charmbracelet/log"
)

var Logger *slog.Logger

func init() {
	l := log.Default()
	sl := slog.New(l)
	Logger = sl
}
