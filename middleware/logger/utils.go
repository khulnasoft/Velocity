package logger

import (
	"io"

	"github.com/khulnasoft/velocity"
	velocitylog "github.com/khulnasoft/velocity/log"
	"github.com/khulnasoft/velocity/utils"
)

func methodColor(method string, colors velocity.Colors) string {
	switch method {
	case velocity.MethodGet:
		return colors.Cyan
	case velocity.MethodPost:
		return colors.Green
	case velocity.MethodPut:
		return colors.Yellow
	case velocity.MethodDelete:
		return colors.Red
	case velocity.MethodPatch:
		return colors.White
	case velocity.MethodHead:
		return colors.Magenta
	case velocity.MethodOptions:
		return colors.Blue
	default:
		return colors.Reset
	}
}

func statusColor(code int, colors velocity.Colors) string {
	switch {
	case code >= velocity.StatusOK && code < velocity.StatusMultipleChoices:
		return colors.Green
	case code >= velocity.StatusMultipleChoices && code < velocity.StatusBadRequest:
		return colors.Blue
	case code >= velocity.StatusBadRequest && code < velocity.StatusInternalServerError:
		return colors.Yellow
	default:
		return colors.Red
	}
}

type customLoggerWriter struct {
	loggerInstance velocitylog.AllLogger
	level          velocitylog.Level
}

func (cl *customLoggerWriter) Write(p []byte) (int, error) {
	switch cl.level {
	case velocitylog.LevelTrace:
		cl.loggerInstance.Trace(utils.UnsafeString(p))
	case velocitylog.LevelDebug:
		cl.loggerInstance.Debug(utils.UnsafeString(p))
	case velocitylog.LevelInfo:
		cl.loggerInstance.Info(utils.UnsafeString(p))
	case velocitylog.LevelWarn:
		cl.loggerInstance.Warn(utils.UnsafeString(p))
	case velocitylog.LevelError:
		cl.loggerInstance.Error(utils.UnsafeString(p))
	default:
		return 0, nil
	}

	return len(p), nil
}

// LoggerToWriter is a helper function that returns an io.Writer that writes to a custom logger.
// You can integrate 3rd party loggers such as zerolog, logrus, etc. to logger middleware using this function.
//
// Valid levels: velocitylog.LevelInfo, velocitylog.LevelTrace, velocitylog.LevelWarn, velocitylog.LevelDebug, velocitylog.LevelError
func LoggerToWriter(logger velocitylog.AllLogger, level velocitylog.Level) io.Writer {
	// Check if customLogger is nil
	if logger == nil {
		velocitylog.Panic("LoggerToWriter: customLogger must not be nil")
	}

	// Check if level is valid
	if level == velocitylog.LevelFatal || level == velocitylog.LevelPanic {
		velocitylog.Panic("LoggerToWriter: invalid level")
	}

	return &customLoggerWriter{
		level:          level,
		loggerInstance: logger,
	}
}
