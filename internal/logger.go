package internal

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

func NewLogger() *slog.Logger {
	return slog.New(
		log.NewWithOptions(
			os.Stdout,
			log.Options{
				Level:           log.InfoLevel,
				ReportTimestamp: true,
			},
		),
	)
}
