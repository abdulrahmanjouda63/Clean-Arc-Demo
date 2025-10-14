package global

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Redis  *redis.Client
	Logger *zap.Logger
)

// InitLogger initializes zap logger with custom configuration
func InitLogger(level, outputPath string) error {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Console encoder for development
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// File encoder for production
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// Configure output
	consoleOutput := zapcore.Lock(os.Stdout)

	var fileOutput zapcore.WriteSyncer
	if outputPath != "" && outputPath != "stdout" {
		// Create directory if it doesn't exist
		dir := filepath.Dir(outputPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		fileOutput = zapcore.AddSync(file)
	}

	// Create core
	var core zapcore.Core
	if fileOutput != nil {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleOutput, zapLevel),
			zapcore.NewCore(fileEncoder, fileOutput, zapLevel),
		)
	} else {
		core = zapcore.NewCore(consoleEncoder, consoleOutput, zapLevel)
	}

	// Create logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

// SyncLogger flushes any buffered log entries
func SyncLogger() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}
