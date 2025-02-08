package config

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	var err error
	config := zap.NewProductionConfig()

	config.Level.SetLevel(getLogLevel(Cfg.LogLevel))

	Logger, err = config.Build()
	if err != nil {
		log.Fatalf("Erro ao inicialiar logger: %v", err)
	}
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
