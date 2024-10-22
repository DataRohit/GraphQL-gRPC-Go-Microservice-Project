package common

import (
	"log"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitializeLogger(logLevel zapcore.Level, outputPaths []string, logFile string) error {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var core zapcore.Core
	if logFile != "" {
		writer, err := getLogWriter(logFile)
		if err != nil {
			return err
		}
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.AddSync(writer),
			logLevel,
		)
	} else {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.AddSync(os.Stdout),
			logLevel,
		)
	}

	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Failed to sync logger: %v", err)
		}
	}()

	return nil
}

func getLogWriter(logFile string) (zapcore.WriteSyncer, error) {
	rotation, err := rotatelogs.New(
		logFile+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(logFile),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		return nil, err
	}
	return zapcore.AddSync(rotation), nil
}

func GetLogger() *zap.Logger {
	if logger == nil {
		err := InitializeLogger(zapcore.InfoLevel, []string{"stdout"}, "")
		if err != nil {
			log.Fatalf("Failed to initialize logger: %v", err)
		}
	}
	return logger
}
