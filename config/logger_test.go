package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestInitLogger(t *testing.T) {
	// Simular nível de log
	Cfg.LogLevel = "debug"

	// Inicializar logger
	InitLogger()
	defer Logger.Sync()

	// Verificar se o logger foi inicializado corretamente
	assert.NotNil(t, Logger)
}

func TestGetLogLevel(t *testing.T) {
	// Testando diferentes níveis de log
	assert.Equal(t, zapcore.DebugLevel, getLogLevel("debug"))
	assert.Equal(t, zapcore.InfoLevel, getLogLevel("info"))
	assert.Equal(t, zapcore.WarnLevel, getLogLevel("warn"))
	assert.Equal(t, zapcore.ErrorLevel, getLogLevel("error"))
	assert.Equal(t, zapcore.InfoLevel, getLogLevel("invalid")) // Deve retornar info como padrão
}
