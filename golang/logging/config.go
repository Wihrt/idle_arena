package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetConfig() zap.Config {
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		DisableStacktrace: false,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg
}
