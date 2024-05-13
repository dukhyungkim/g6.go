package lib

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type AppLogger struct {
	zerolog.Logger
}

var Logger AppLogger

func NewLogger() AppLogger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	// format error
	output.FormatErrFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s: ", i)
	}

	logger := zerolog.New(output).With().Caller().Timestamp().Logger()
	Logger = AppLogger{logger}
	return Logger
}

func (l *AppLogger) LogInfo() *zerolog.Event {
	return l.Logger.Info()
}

func (l *AppLogger) LogError() *zerolog.Event {
	return l.Logger.Error()
}

func (l *AppLogger) LogDebug() *zerolog.Event {
	return l.Logger.Debug()
}

func (l *AppLogger) LogWarn() *zerolog.Event {
	return l.Logger.Warn()
}

func (l *AppLogger) LogFatal() *zerolog.Event {
	return l.Logger.Fatal()
}
