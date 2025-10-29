package gcp

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewGCPLogger() (*zap.Logger, error) {
	loggerCfg := &zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return loggerCfg.Build(zap.AddStacktrace(zap.DPanicLevel))
}

var encoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "severity",
	NameKey:        "logger-presets",
	CallerKey:      "caller",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    encodeLevel(),
	EncodeTime:     zapcore.RFC3339TimeEncoder,
	EncodeDuration: zapcore.MillisDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

func encodeLevel() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString("DEBUG")
		case zapcore.InfoLevel:
			enc.AppendString("INFO")
		case zapcore.WarnLevel:
			enc.AppendString("WARNING")
		case zapcore.ErrorLevel:
			enc.AppendString("ERROR")
		case zapcore.DPanicLevel:
			enc.AppendString("CRITICAL")
		case zapcore.PanicLevel:
			enc.AppendString("ALERT")
		case zapcore.FatalLevel:
			enc.AppendString("EMERGENCY")
		default:
			enc.AppendString("DEFAULT")
		}
	}
}
