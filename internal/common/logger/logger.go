// Package logger represents logger
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New return new logger instance
func New() (*zap.SugaredLogger, error) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.DisableCaller = true
	zapConfig.DisableStacktrace = true
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zapLogger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return zapLogger.Sugar(), nil
}
