package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init creates a production-ready logger based on the log level provided.
func Init() (*zap.Logger, error) {
	level := os.Getenv("LOG_LEVEL")
	if len(level) == 0 {
		level = "info"
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(toZapLevel(level))

	// Make timestamps human-readable instead of Unix timestamps
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return config.Build()
}

func toZapLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}
