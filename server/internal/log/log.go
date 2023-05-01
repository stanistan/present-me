package log

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func Ctx(ctx context.Context) *zerolog.Logger {
	return log.Ctx(ctx)
}

func init() {
	// This is sset to be GOOGLE format (ish)
	// - https://cloud.google.com/logging/docs/structured-logging
	// - https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "times"
	zerolog.MessageFieldName = "message"
	zerolog.ErrorFieldName = "message"
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		switch l {
		case zerolog.TraceLevel:
			return "DEBUG"
		case zerolog.DebugLevel:
			return "DEBUG"
		case zerolog.InfoLevel:
			return "INFO"
		case zerolog.WarnLevel:
			return "WARNING"
		case zerolog.ErrorLevel:
			return "ERROR"
		case zerolog.FatalLevel:
			return "CRITICAL"
		case zerolog.PanicLevel:
			return "ALERT"
		case zerolog.NoLevel:
			return "DEFAULT"
		}
		return "DEFAULT"
	}
}
