package logger

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"os"
	"time"
)

// TODO: возвращать slog, чтобы не зависить от конкретной реализации стороннего логгера
func New() *log.Logger {
	const MaxWidth = 7
	const Black = "#000000"
	const White = "#ffffff"
	const Green = "#28a745"
	const Yellow = "#ffc107"
	const Red = "#dc3545"

	options := log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
	}

	logger := log.NewWithOptions(os.Stdout, options)

	styles := log.DefaultStyles()

	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Bold(true).
		MaxWidth(MaxWidth).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color(Green)).
		Foreground(lipgloss.Color(White))

	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("WARN").
		Bold(true).
		MaxWidth(MaxWidth).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color(Yellow)).
		Foreground(lipgloss.Color(Black))

	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERROR").
		Bold(true).
		MaxWidth(MaxWidth).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color(Red)).
		Foreground(lipgloss.Color(White))

	logger.SetStyles(styles)

	return logger
}
